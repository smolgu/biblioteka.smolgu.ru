package vktracker

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/smolgu/lib/modules/setting"
	dry "github.com/ungerik/go-dry"
	"github.com/zhuharev/intarr"
	"github.com/zhuharev/vkutil"
)

var nameCache = map[int]vkutil.User{}
var api *vkutil.Api

func NewContext() {
	go func() {
		for {
			job()
			time.Sleep(6 * time.Hour)
		}
	}()
}

func Reports(limit int) ([]LogRecord, error) {
	return readLog(limit)
}

func job() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	api = vkutil.New()
	api.SetDebug(false)
	api.VkApi.Lang = "ru"
	api.VkApi.AccessToken = setting.VkAccessToken
	//err := api.DirectAuth("unybigk@gmail.com", "meepohazhu")
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Println(api.VkApi.AccessToken)

	members, err := api.GroupsGetAllMembers(setting.VkGroupID)
	if err != nil {
		log.Println(err)
		return
	}
	oldMembers, err := readIds()
	if err != nil {
		log.Println(err)
		return
	}

	users, err := api.UsersGet(members, url.Values{
		"fields": {"photo"},
	})
	if err != nil {
		log.Println(err)
		return
	}

	for _, v := range users {
		nameCache[v.Id] = v
	}

	log.Println(users)

	j, r := intarr.Diff(oldMembers, members)
	if len(j) > 0 || len(r) > 0 {
		err = writeLog(j, r)
		if err != nil {
			log.Println(err)
		}
	}

	err = dry.FileSetJSON("data/ids", members)
	if err != nil {
		log.Println(err)
	}

	off = 0
}

type LogRecord struct {
	Joined []int
	Leaved []int
	Date   time.Time
	Users  map[int]vkutil.User
}

func (l LogRecord) Get(id int) vkutil.User {
	return l.Users[id]
}

func (l LogRecord) Img(id int) string {
	return l.Users[id].Photo
}

func (l LogRecord) Fname(id int) string {
	return l.Users[id].FirstName
}

func (l LogRecord) Lname(id int) string {
	return l.Users[id].LastName
}

var (
	off int64
)

func revRead(f *os.File) (uint64, error) {

	// _, err = f.Seek(info.Size(), io.SeekStart)
	// if err != nil {
	// 	return 0, err
	// }
	var res uint64
	err := binary.Read(f, binary.BigEndian, &res)
	if err != nil {
		return 0, err
	}

	if off-16 < 0 {
		off = -1
		return res, nil
	}
	off, err = f.Seek(-16, io.SeekCurrent)
	if err != nil {
		log.Println("rev", err)
		return 0, err
	}
	return res, nil
}

func readLog(limit int) ([]LogRecord, error) {
	var (
		res       []LogRecord
		curRecNum = 1
	)

	f, err := os.OpenFile("data/vk.log", os.O_CREATE|os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	off, err = f.Seek(-8, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	for curRecNum <= limit && off != -1 {
		curRec := LogRecord{}
		i, err := revRead(f)
		if err != nil {
			return nil, err
		}
		t := time.Unix(0, int64(i))
		curRec.Date = t

		leavedLen, err := revRead(f)
		if err != nil {
			return nil, err
		}

		joinedLen, err := revRead(f)
		if err != nil {
			return nil, err
		}

		curRec.Joined = make([]int, int(joinedLen))
		curRec.Leaved = make([]int, int(leavedLen))

		for i := range curRec.Leaved {
			lolka, err := revRead(f)
			if err != nil {
				return nil, err
			}
			get(int(lolka))
			curRec.Leaved[i] = int(lolka)
		}

		for i := range curRec.Joined {
			lolka, err := revRead(f)
			if err != nil {
				return nil, err
			}
			get(int(lolka))
			curRec.Joined[i] = int(lolka)
		}
		log.Println(curRec)
		curRec.Users = nameCache
		res = append(res, curRec)
		curRecNum++
	}

	return res, nil
}

func get(id int) (vkutil.User, error) {
	if user, has := nameCache[id]; has {
		return user, nil
	}
	users, err := api.UsersGet(id, url.Values{
		"fields": {"photo"},
	})
	if err != nil {
		return vkutil.User{}, err
	}
	if len(users) != 1 {
		return vkutil.User{}, fmt.Errorf("lol")
	}
	nameCache[id] = users[0]
	return users[0], nil
}

func writeLog(joined, leaved []int) error {
	f, err := os.OpenFile("data/vk.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	bw := bufio.NewWriter(f)
	for _, v := range joined {
		err = binary.Write(bw, binary.BigEndian, uint64(v))
		if err != nil {
			return err
		}
	}
	for _, v := range leaved {
		err = binary.Write(bw, binary.BigEndian, uint64(v))
		if err != nil {
			return err
		}
	}
	err = binary.Write(bw, binary.BigEndian, uint64(len(joined)))
	if err != nil {
		return err
	}
	err = binary.Write(bw, binary.BigEndian, uint64(len(leaved)))
	if err != nil {
		return err
	}
	err = binary.Write(bw, binary.BigEndian, uint64(time.Now().UnixNano()))
	if err != nil {
		return err
	}

	err = bw.Flush()
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

func readIds() ([]int, error) {
	if !dry.FileExists("data/ids") {
		return nil, nil
	}
	bts, err := ioutil.ReadFile("data/ids")
	if err != nil {
		return nil, err
	}
	var ids []int
	err = json.Unmarshal(bts, &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

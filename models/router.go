package models

type ObjectType int

const (
	ObjectPage ObjectType = iota + 1
)

type Router struct {
	Slug     string `xorm:"unique index"`
	Type     ObjectType
	ObjectID int64
}

func RouterSave(objectID int64, objectType ObjectType, slug string) error {
	r := new(Router)
	r.ObjectID = objectID
	r.Slug = slug
	r.Type = objectType
	_, err := x.Insert(r)
	return err
}

func RouterResolve(slug string) (objectID int64, typ ObjectType, err error) {
	r := new(Router)
	has, err := x.Where("slug = ?", slug).Get(r)
	if err != nil {
		return
	}
	if !has {
		err = ErrNotFound
		return
	}
	return r.ObjectID, r.Type, nil
}

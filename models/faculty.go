package models

type Faculty struct {
	Id    int64
	Title string
	Short string
}

type Faculties []Faculty

func (f Faculties) Slice() (result []string) {
	for _, v := range f {
		result = append(result, v.Title)
	}
	return
}

var SmolGUFaculties = Faculties{
	{1, "Физико-математический факультет", "ФМФ"},
	{2, "Филологический факультет", "ФМФ"},
	{3, "Факультет истории и права", "ФИиП"},
	{4, "Естественно-географический факультет", "ЕГФ"},
	{5, "Социальный факультет", "СФ"},
	{6, "Психолого-педагогический факультет", "ППФ"},
	{7, "Художественно-графический факультет", "ППФ"},
	{8, "Факультет экономики и управления", "ФЭУ"},
	{9, "Другое", "Др"},
}

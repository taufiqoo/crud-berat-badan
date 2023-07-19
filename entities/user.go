package entities

type User struct {
	Id      int64
	Tanggal string `validate:"required"`
	Max     string `validate:"required" label:"Berat badan maksimal"`
	Min     string `validate:"required" label:"Berat badan minimal"`
	Diff    int
}

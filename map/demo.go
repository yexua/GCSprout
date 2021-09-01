package main

var list map[string]*Student

type People interface {
	Speak(string) string
}

type Student struct {
	Name string
}

func (stu *Student) Speak(think string) (talk string) {
	if think == "love" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return
}

func main() {
	list = make(map[string]*Student)

	student := Student{"Jack"}

	list["student"] = &student
	list["student"].Name = "Tom"

}

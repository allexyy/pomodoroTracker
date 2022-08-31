package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"os"
	"path/filepath"
	"pomodoro/window"
	"time"
)

type Task struct {
	Name        string
	Timer       time.Duration
	RepeatCount int
}

var data = setTaskListData()
var myApp = app.New()
var taskChan = make(chan bool, 1)
var taskIterationTime = 10 * time.Second

func main() {
	go func() { taskChan <- false }()
	myWindow := myApp.NewWindow("Table Widget")
	myWindow.Resize(fyne.NewSize(800, 400))

	mainContent := window.GetContent()
	settingPage := window.GetSettingsPage()

	mainContent.Add(widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), func() {
		myWindow.SetContent(settingPage)
	}))
	taskList := container.New(layout.NewVBoxLayout())

	input := widget.NewEntry()
	input.Resize(fyne.NewSize(150, 50))
	input.SetPlaceHolder("Text task Name")
	b := widget.NewButton("Add Task", func() {
		if input.Text != "" {
			newTask := Task{Name: input.Text, Timer: taskIterationTime, RepeatCount: 0}
			data[input.Text] = newTask
			updateTaskList(data, taskList)
			saveTask(data)
		}
		input.SetText("")
	})

	addTaskContent := container.NewVBox(input, b)
	addTaskContent.Resize(fyne.NewSize(150, 50))

	updateTaskList(data, taskList)

	taskPage := container.New(layout.NewVBoxLayout(), mainContent, addTaskContent, taskList)

	myWindow.SetContent(taskPage)

	backToMainPage := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		myWindow.SetContent(taskPage)
	})
	backToMainPage.Resize(fyne.NewSize(123, 300))
	settingPage.Add(backToMainPage)

	myWindow.ShowAndRun()
}

func updateTaskList(data map[string]Task, taskList *fyne.Container) {
	taskList.RemoveAll()
	for _, taskFromList := range data {
		currentTaskRepeat := string(taskFromList.RepeatCount)
		currentTaskTime := fmt.Sprintf("%f", taskFromList.Timer.Minutes())
		timeText := canvas.NewText(currentTaskTime, color.White)
		startIcon := theme.MediaPlayIcon()
		taskContainer := container.New(
			layout.NewHBoxLayout(),
			canvas.NewText(taskFromList.Name, color.White),
			widget.NewButtonWithIcon("", startIcon, func() {
				go taskFromList.startTask(timeText, startIcon)
				taskList.Refresh()

			}),
			widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
				taskFromList.deleteTask(taskList)
			}),
			canvas.NewText(currentTaskRepeat, color.White),
			timeText)

		taskList.Add(taskContainer)
	}
}

func (currentTask *Task) startTask(text *canvas.Text, icon fyne.Resource) {
	icon = theme.MediaStopIcon()
	isStart := <-taskChan
	if isStart == true {
		icon = theme.MediaPlayIcon()
		taskChan <- false
		return
	}
	updateTime(text, currentTask.Timer)
	taskChan <- true
	go func() {
		for range time.Tick(time.Second) {
			fmt.Println("Refresh time")
			currentTask.Timer = updateTime(text, currentTask.Timer)
			text.Refresh()
			if currentTask.Timer.Seconds() == 0 {
				fmt.Println("I DO")
				currentTask.endTaskTimer()
				text.Refresh()
				fmt.Println(currentTask.Timer)
				fmt.Println("Refresh")
				return
			}
		}
	}()

}

func updateTime(text *canvas.Text, taskTime time.Duration) time.Duration {
	taskTime = taskTime - 1*time.Second
	text.Text = taskTime.String()
	return taskTime
}

func (currentTask *Task) deleteTask(taskList *fyne.Container) {
	delete(data, currentTask.Name)
	updateTaskList(data, taskList)
}

func (currentTask *Task) endTaskTimer() {
	currentTask.RepeatCount++
	currentTask.Timer = taskIterationTime
}

func isFileExist(filename string) bool {
	_, err := os.Open(filename)
	result := true

	if err != nil {
		result = false
	}

	fmt.Println(result)
	return result
}

func writeDataToFile(file *os.File, dataToWrite []byte) {
	fmt.Println(string(dataToWrite))
	result, err := file.WriteString(string(dataToWrite))

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Write result")
	fmt.Println(result)

}

func saveTask(taskList map[string]Task) {
	fmt.Println(taskList)
	jsonString, err := json.Marshal(taskList)
	fmt.Println(jsonString)

	if err != nil {
		fmt.Println(err)
	}

	if isFileExist("task-list.json") == false {
		fmt.Println("Create File")
		file, err := os.Create(filepath.Join("Go", "src", "pomodoroTrackeer", "task-list.json"))
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		}
		writeDataToFile(file, jsonString)
	} else {
		fmt.Println("Open File")
		file, err := os.Open("task-list.json")
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		}
		writeDataToFile(file, jsonString)
	}

}

func setTaskListData() map[string]Task {
	//TODO: Load data from json
	return make(map[string]Task)
}

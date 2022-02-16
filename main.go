package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	frame_width   = 480
	frame_height  = 480
	Width_count   = 12
	Height_count  = 16
	Height        = 30
	Width         = 40
	default_dir   = "img/"
	res_dir       = "res/"
	final_res_dir = "final_res/"
	letter_dir    = "letters/"
)

type Changeable interface {
	Set(x, y int, c color.Color)
}

type cell struct {
	value               string
	not_accepted_values []string
}

func new_cell(value string, not_accepted_values []string) cell {
	return cell{
		value:               value,
		not_accepted_values: not_accepted_values,
	}
}

func roll_rand(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

func def_max_flag_count(flag_count int, grid_size int) (int, []int) {
	max := make([]int, flag_count)
	for i := range max {
		max[i] = 0
	}
	return int(grid_size/flag_count) + 1, max
}
func print_grid(grid [][]cell) {
	for i := 0; i < Height_count; i++ {
		for j := 0; j < Width_count; j++ {
			fmt.Print(grid[i][j].value, " ")
		}
		fmt.Println()
	}
}
func create_rand_grid(grid [][]cell, flag []string) {
	flags_max, stat := def_max_flag_count(len(flag), Width_count*Height_count)
	for i := 0; i < Height_count; i++ {
		for j := 0; j < Width_count; j++ {
			var x int
			//ac accepted value
			//ac accepted count
			for av, ac := true, false; !(av && ac); {
				av, ac = true, false
				x = roll_rand(len(flag))
				for _, v := range grid[i][j].not_accepted_values {
					if v == flag[x] {
						av = false
					}
				}

				if stat[x] < flags_max {
					ac = true

				}
			}
			grid[i][j].value = flag[x]
			stat[x]++
			if i > 0 {
				grid[i-1][j].not_accepted_values = append(grid[i-1][j].not_accepted_values, flag[x])
			}
			if i < Height_count-1 {
				grid[i+1][j].not_accepted_values = append(grid[i+1][j].not_accepted_values, flag[x])
			}
			if j > 0 {
				grid[i][j-1].not_accepted_values = append(grid[i][j-1].not_accepted_values, flag[x])
			}
			if j < Width_count-1 {
				grid[i][j+1].not_accepted_values = append(grid[i][j+1].not_accepted_values, flag[x])
			}

		}

	}
}
func save_img(path string, img image.Image) error {
	f, _ := os.Create(path + ".png")
	err := png.Encode(f, img)
	defer f.Close()
	return err
}
func create_rand_bg(grid [][]cell, m map[string]image.Image, id int) {
	for j := 0; j < Height_count; j++ {
		for i := 0; i < Width_count; i++ {
			for y := 0; y < Height; y++ {
				for x := 0; x < Width; x++ {
					m["template"].(Changeable).Set(i*Width+x, j*Height+y, m[grid[j][i].value].At(x, y))
				}
			}
		}
	}
	save_img(res_dir+strconv.Itoa(id), m["template"])
}
func walk_dir(path string) map[string]image.Image {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	m := make(map[string]image.Image)

	for _, file := range files {
		img := decode_png(path, file.Name())
		m[file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]] = img
	}
	return m
}
func reset_grid(grid [][]cell) {
	for i := 0; i < Height_count; i++ {
		grid[i] = make([]cell, Width_count)
		for j := 0; j < Width_count; j++ {
			grid[i][j] = new_cell("", []string{})
		}
	}
}

func fit_letter(letters map[string]image.Image, m map[string]image.Image) {
	for lk := range letters {
		for ik := range m {
			letter_copy := decode_png(letter_dir, lk+".png")
			for x := 0; x <= frame_width; x++ {
				for y := 0; y <= frame_height; y++ {
					r, _, _, _ := letter_copy.At(x, y).RGBA()
					if r != 0 {
						letter_copy.(Changeable).Set(x, y, m[ik].At(x, y))
					}

				}
			}
			save_img(final_res_dir+lk+ik, letter_copy)
		}

	}

}
func decode_png(dir string, name string) image.Image {

	imgfile, err := os.Open(dir + name)
	if err != nil {
		panic(err.Error())
	}
	defer imgfile.Close()
	img, err := png.Decode(imgfile)
	if err != nil {
		panic(err.Error())
	}
	return img
}
func main() {

	m := walk_dir(default_dir)

	grid := make([][]cell, Height_count)
	reset_grid(grid)
	flags_with_i := []string{
		"in",
		"id",
		"ir",
		"it",
		"iq",
		"il",
		"ie",
		"is",
	}
	for k := 0; k < len(flags_with_i); k++ {
		reset_grid(grid)
		create_rand_grid(grid, flags_with_i)
		fmt.Printf("grid #%d\n", k)
		print_grid(grid)
		create_rand_bg(grid, m, k)
	}
	letters := walk_dir(letter_dir)
	res := walk_dir(res_dir)
	fit_letter(letters, res)

}

package main

import (
	"flag"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"sort"
)

//原图 image.Image
//↓ getPixelData
//二维灰度切片 [][]uint8
//↓ makeImmutableMatrix
//只读闭包 func(y, x int) uint8
//↓ medianFilter(0, height, 0, width, ...)
//滤波后的二维切片 [][]uint8
//↓ flattenImage
//一维切片 []uint8
//↓ image.NewGray + png.Encode
//输出图片 out.png

// check handles a potential error.
// It stops execution of the program ("panics") if an error has happened.
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// makeMatrix makes and returns a 2D slice with the given dimensions.
func makeMatrix(height, width int) [][]uint8 {
	matrix := make([][]uint8, height)
	for i := range matrix {
		matrix[i] = make([]uint8, width)
	}
	return matrix
}

// makeImmutableMatrix takes an existing 2D matrix and wraps it in a getter closure.
func makeImmutableMatrix(matrix [][]uint8) func(y, x int) uint8 {
	return func(y, x int) uint8 {
		return matrix[y][x]
	}
}

// medianFilter applies the filter between the given x and y bounds on the given closure.
// medianFilter returns the section where the filter was applied as a 2D slice.
func medianFilter(startY, endY, startX, endX int, data func(y, x int) uint8) [][]uint8 {
	height := endY - startY
	width := endX - startX
	radius := 2               //滤波半径
	midPoint := (5*5 + 1) / 2 //中点

	filteredMatrix := makeMatrix(height, width) //初始化
	filterValues := make([]int, 5*5)

	for i := radius + startY; i < endY-radius; i++ {
		for j := radius + startX; j < endX-radius; j++ {
			count := 0
			for k := i - radius; k <= i+radius; k++ {
				for l := j - radius; l <= j+radius; l++ {
					filterValues[count] = int(data(k, l))
					count++
				}
			}
			sort.Ints(filterValues)
			filteredMatrix[i-startY][j-startX] = uint8(filterValues[midPoint]) //对值排序
		}
	}
	return filteredMatrix
}

// getPixelData transfers an image.Image to a standard 2D slice.
func getPixelData(img image.Image) [][]uint8 {
	bounds := img.Bounds()
	pixels := makeMatrix(bounds.Dy(), bounds.Dx())

	curr := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			pixels[y][x] = uint8(lum / 256)
			curr++
		}
	}
	return pixels
}

// loadImage opens a file and returns the contents as an image.Image.
func loadImage(filepath string) image.Image {
	existingImageFile, err := os.Open(filepath)
	check(err)
	defer existingImageFile.Close()

	img, _, err := image.Decode(existingImageFile)
	check(err)

	return img
}

// flattenImage takes a 2D slice and flattens it into a single 1D slice.
func flattenImage(image [][]uint8) []uint8 {
	height := len(image)
	width := len(image[0])

	flattenedImage := make([]uint8, 0, height*width)
	for i := 0; i < height; i++ {
		flattenedImage = append(flattenedImage, image[i]...)
	}
	return flattenedImage
}

// filter reads in a png image, applies the filter and outputs the result as a png image.
// filter is the function called by the tests in medianfilter_test.go
func filter(filepathIn, filepathOut string, threads int) {
	image.RegisterFormat("png", "PNG", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

	img := loadImage(filepathIn)
	bounds := img.Bounds()
	height := bounds.Dy()
	width := bounds.Dx()

	immutableData := makeImmutableMatrix(getPixelData(img)) //getPixelData(img) 把整张图转成 [][]uint8 灰度矩阵
	var newPixelData [][]uint8

	//并行化的目标就是：把 medianFilter 的大任务拆成若干小块，交给多个 goroutine 同时算，再拼回整张图
	if threads < 1 {
		threads = 1
	}

	if threads == 1 { //单线程保持原样 //thread 就是几个worker一起干活
		newPixelData = medianFilter(0, height, 0, width, immutableData) //处理整张图，从第 0 行到第 height-1 行，从第 0 列到第 width-1 列
	} else {
		newPixelData = make([][]uint8, 0, height) //创建一个空的二维切片，但预留 height 行的容量

		// 1. 顶部边界：最上面 2 行保持黑色
		for y := 0; y < 2 && y < height; y++ {
			newPixelData = append(newPixelData, make([]uint8, width))
		}

		// 2. 计算内部有效行，分块，启动 worker
		innerHeight := height - 4
		if innerHeight > 0 {
			chunkSize := (innerHeight + threads - 1) / threads //平均一个threads 工作至少多少

			//make(类型, 长度, 容量)
			//第一个【】是slide//chan [][]uint8：切片的每个元素是一个通道
			outChans := make([]chan [][]uint8, 0, threads)

			for t := 0; t < threads; t++ { //t 是第几个worker
				outStartY := 2 + t*chunkSize
				if outStartY >= height-2 { // 注意：>=，不是 >
					break
				}

				outEndY := 2 + (t+1)*chunkSize //output 结束行
				if outEndY > height-2 {
					outEndY = height - 2
				}

				startY := outStartY - 2 // 往上多拿 2 行
				endY := outEndY + 2     // 往下多拿 2 行

				out := make(chan [][]uint8, 1)   //给这个 worker 用的通道，缓冲区大小 1
				outChans = append(outChans, out) //新切片 = append(原切片, 新元素...)把通道保存起来，后面按顺序收结果

				go worker(startY, endY, 0, width, immutableData, out) //启动一个新的 goroutine，并行执行 worker 函数。
			}

			// 3. 按顺序收集结果，只取有效行，append 到 newPixelData
			for _, out := range outChans {
				result := <-out
				// result 的前 2 行和最后 2 行是额外邻域，不作为输出
				for y := 2; y < len(result)-2; y++ {
					newPixelData = append(newPixelData, result[y])
				}
			}
		}

		// 4. 底部边界：最下面 2 行保持黑色
		for y := height - 2; y < height; y++ {
			if y >= 0 {
				newPixelData = append(newPixelData, make([]uint8, width))
			}
		}
	}

	imout := image.NewGray(image.Rect(0, 0, width, height))
	imout.Pix = flattenImage(newPixelData) //返回的 newPixelData 就是整张图滤波后的矩阵
	ofp, _ := os.Create(filepathOut)
	defer ofp.Close()
	err := png.Encode(ofp, imout)
	check(err)
}

// Q1B
func worker(startY, endY, startX, endX int, data func(y, x int) uint8, out chan<- [][]uint8) {
	// 一行：调用 medianFilter
	out <- medianFilter(startY, endY, startX, endX, data)
	// 一行：把结果发送到 out
}

// main reads in the filepath flags or sets them to default values and calls filter().
func main() {
	var filepathIn string
	var filepathOut string
	var threads int

	flag.StringVar(
		&filepathIn,
		"in",
		"ship.png",
		"Specify the input file.")

	flag.StringVar(
		&filepathOut,
		"out",
		"out.png",
		"Specify the output file.")

	flag.IntVar(
		&threads,
		"threads",
		1,
		"Specify the number of worker threads to use.")

	flag.Parse()
	filter(filepathIn, filepathOut, threads)
}

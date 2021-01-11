package internal

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

func SaveStatistics(filename string, timings [][]time.Duration, fileSizes []FileSize) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	fileSizesLen := len(fileSizes)

	means := make([]time.Duration, fileSizesLen)
	for _, timing := range timings {
		for sizeIndex := 0; sizeIndex < fileSizesLen; sizeIndex++ {
			means[sizeIndex] += timing[sizeIndex]
		}
	}

	for sizeIndex := 0; sizeIndex < fileSizesLen; sizeIndex++ {
		means[sizeIndex] /= 5
	}

	stdDeviationsFloats := make([]float64, fileSizesLen)
	for _, timing := range timings {
		for sizeIndex := 0; sizeIndex < fileSizesLen; sizeIndex++ {
			minusMean := timing[sizeIndex].Seconds() - means[sizeIndex].Seconds()
			stdDeviationsFloats[sizeIndex] += minusMean * minusMean
		}
	}

	for sizeIndex := 0; sizeIndex < fileSizesLen; sizeIndex++ {
		stdDeviationsFloats[sizeIndex] = math.Sqrt(stdDeviationsFloats[sizeIndex] / 4.0)
	}

	_, err = file.WriteString("File Size (bytes),Mean Elapsed Time\n")
	if err != nil {
		log.Fatal(err)
	}
	for sizeIndex := 0; sizeIndex < fileSizesLen; sizeIndex++ {
		_, err := file.WriteString(fmt.Sprintf("%v,\"%v ± %vs\"\n", fileSizes[sizeIndex], means[sizeIndex], stdDeviationsFloats[sizeIndex]))
		if err != nil {
			log.Fatal(err)
		}
	}
}

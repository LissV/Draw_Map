package main

import (
	"fmt"
	"io"
	"os"

	"github.com/fogleman/gg"
	geojson "github.com/paulmach/go.geojson"
)

func main() {
	file, err := os.Open("regions.geojson")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	rawFeatureJSON := make([]byte, 0)
	dataSclice := make([]byte, 100)

	for {
		n, err := file.Read(dataSclice)
		if err == io.EOF {
			break
		}
		rawFeatureJSON = append(rawFeatureJSON, dataSclice[:n]...)
	}

	fc, err := geojson.UnmarshalFeatureCollection(rawFeatureJSON)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	dc := gg.NewContext(1366, 1024)
	dc.SetHexColor("fff")

	dc.Scale(7, 7)

	for i := 0; i < len(fc.Features); i++ {
		for j := 0; j < len(fc.Features[i].Geometry.MultiPolygon); j++ {
			for k := 0; k < len(fc.Features[i].Geometry.MultiPolygon[j]); k++ {
				for m := 0; m < len(fc.Features[i].Geometry.MultiPolygon[j][k])-1; m++ {
					x1 := float64(fc.Features[i].Geometry.MultiPolygon[j][k][m][0])
					y1 := float64(fc.Features[i].Geometry.MultiPolygon[j][k][m][1])
					x2 := float64(fc.Features[i].Geometry.MultiPolygon[j][k][m+1][0])
					y2 := float64(fc.Features[i].Geometry.MultiPolygon[j][k][m+1][1])
					if m == 0 {
						dc.DrawLine(x1, y1, x2, y2)
					} else {
						dc.LineTo(x2, y2)
					}
				}
			}
		}
	}

	dc.SetHexColor("f0f")
	dc.Fill()
	dc.SavePNG("out.png")

}

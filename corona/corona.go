package corona

import (
	"fmt"
	"github.com/dghubble/sling"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"log"
	"net/http"
	"sort"
	"time"
)

type Corona struct {
	Latest    Latest      `json:"latest"`
	Locations []Locations `json:"locations"`
}
type Latest struct {
	Confirmed int `json:"confirmed"`
	Deaths    int `json:"deaths"`
	Recovered int `json:"recovered"`
}
type Coordinates struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
type Timeline map[time.Time]int

type Confirmed struct {
	Latest   int      `json:"latest"`
	Timeline Timeline `json:"timeline"`
}
type Deaths struct {
	Latest   int      `json:"latest"`
	Timeline Timeline `json:"timeline"`
}
type Recovered struct {
	Latest   int      `json:"latest"`
	Timeline Timeline `json:"timeline"`
}
type Timelines struct {
	Confirmed Confirmed `json:"confirmed"`
	Deaths    Deaths    `json:"deaths"`
	Recovered Recovered `json:"recovered"`
}
type Locations struct {
	Coordinates Coordinates `json:"coordinates"`
	Country     string      `json:"country"`
	CountryCode string      `json:"country_code"`
	ID          int         `json:"id"`
	LastUpdated time.Time   `json:"last_updated"`
	Latest      Latest      `json:"latest"`
	Province    string      `json:"province"`
	Timelines   Timelines   `json:"timelines"`
}

type TimeLineData struct {
	Time  time.Time
	Count int
}
type timelineArray []TimeLineData

type CoronaQueryData struct {
	CountryCode string `url:"country_code"`
	Timelines   bool   `url:"timelines"`
}

func GetCoronaData() *Corona {
	var t Corona

	coronaBase := sling.New().Base("https://coronavirus-tracker-api.herokuapp.com/").Client(http.DefaultClient)
	params := &CoronaQueryData{
		CountryCode: "DE",
		Timelines:   true,
	}
	resp, err := coronaBase.New().Get("v2/locations").QueryStruct(params).ReceiveSuccess(&t)
	//log.Printf("%+#v Response %#v, err %v", t, resp, err)
	if err != nil {
		log.Print(err)
	}
	if resp.StatusCode == 200 {
		return &t
	}
	return nil
}

func SaveCoronaImage(c *Corona, image string) {
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}

	xticks := plot.TimeTicks{Format: "02.02.2006"}

	//	log.Printf("Confirmed Timelines %+#v", c.Locations[0].Timelines.Confirmed.Timeline)

	p.Title.Text = fmt.Sprintf("Corona in %s from %s\n\n", c.Locations[0].Country, c.Locations[0].LastUpdated.Format("02.02.2006 15:04:05"))
	p.X.Label.Text = "Datum"
	p.Y.Label.Text = "Personen"
	p.X.Tick.Marker = xticks
	p.Legend.Left = true
	p.Legend.Top = true

	xyConfirmed := GetPlotterXYbyTimeline(c.Locations[0].Timelines.Confirmed.Timeline)
	xyDeaths := GetPlotterXYbyTimeline(c.Locations[0].Timelines.Deaths.Timeline)
	xyRecoverd := GetPlotterXYbyTimeline(c.Locations[0].Timelines.Recovered.Timeline)
	err = plotutil.AddLinePoints(p,
		"Confirmed", xyConfirmed,
		"Deaths", xyDeaths,
		"Recoverd", xyRecoverd,
	)
	if err != nil {
		panic(err)
	}
	//log.Printf("%+#v", xyConfirmed)
	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 3*vg.Inch, image); err != nil {
		log.Fatal(err)
	}

	//img := image.NewRGBA(image.Rect(0, 0, int(4*vg.Inch), int(3*vg.Inch)))
	//d := vgimg.NewWith(vgimg.UseImage(img))
	//p.Draw(draw.New(d))
	//
	//// Save the image.
	//f, err := os.Create("test.png")
	//if err != nil {
	//	panic(err)
	//}
	//if err := png.Encode(f, img); err != nil {
	//	panic(err)
	//}
	//if err := f.Close(); err != nil {
	//	panic(err)
	//}
}

func GetPlotterXYbyTimeline(t Timeline) plotter.XYs {
	sortedConfirmedTimeline := make(timelineArray, 0, len(t))
	for k, v := range t {
		sortedConfirmedTimeline = append(sortedConfirmedTimeline, TimeLineData{
			Time:  k,
			Count: v,
		})
	}
	sort.Slice(sortedConfirmedTimeline, func(i, j int) bool {
		return sortedConfirmedTimeline[i].Time.Before(sortedConfirmedTimeline[j].Time)
	})

	var xy plotter.XYs
	for _, v := range sortedConfirmedTimeline {
		xy = append(xy, plotter.XY{float64(v.Time.Unix()), float64(v.Count)})
	}
	return xy
}

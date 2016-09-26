package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/olekukonko/tablewriter"
	"sync"
	"strings"
	"time"
	"os"
	"strconv"
)

var companies map[int]string = make(map[int]string)

var cities map[int]string = map[int]string{
	1: "Киев",
	2: "Днепр",
	3: "Харьков",
	4: "Одесса",
	5: "Львов",
	6: "Черкассы",
	7: "Винница",
	8: "Николаев",
}

var (
	KEYS string 		= "php,yii"
	CITY_CODE int		= 1
	MAX_POSTION int		= 60
)

func init(){
    flag.StringVar(&KEYS, "k", KEYS, "ключивые слова через запятую")
	flag.IntVar(&CITY_CODE, "cc", CITY_CODE, "id города "+citiesToString())
    flag.IntVar(&MAX_POSTION, "mp", MAX_POSTION, "максимальная позиция в рейтинге")
}

func main() {
	KEYS = strings.ToLower(KEYS)
    findCompanies()
	jobs := grab()
	
	for _, job := range jobs{
		data := [][]string{
			[]string{"Position", job["position"]}, 
			[]string{"Company Name", job["companyName"]},
			[]string{"Job", job["jobName"]},
			[]string{"Description", job["description"]},
			[]string{"Link", job["link"]},
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetBorder(true)
		for _, v := range data {
			table.Append(v)
		}
		table.Render()
		
		/*fmt.Printf("Rate: %s\n", job["position"])
		fmt.Printf("Company: %s\n", job["companyName"])
		fmt.Printf("Job: %s\n", job["jobName"])
		fmt.Printf("Description: %s\n", job["description"])
		fmt.Printf("Link: %s\n\n", job["link"])*/
		
	}
}

func citiesToString() string{
	var total string
	for i, value := range cities {
		total += strconv.Itoa(i)+" - "+value+";"
	}
	return total
}

func findCompanies(){
	page, err := goquery.NewDocument("https://jobs.dou.ua/ratings/%D0%9A%D0%B8%D0%B5%D0%B2/")
	if err != nil {
        panic(err)
    }
	
	page.Find(".company-name a").Each(func(i int, s *goquery.Selection) {
		if i+1 <= MAX_POSTION {
			companies[i+1]=s.Text()
		}		
    })
}

func isKeysContains(keys string, text string)bool{
	text = strings.ToLower(text)
	stringSlice := strings.Split(keys, ",")
	for _, key := range stringSlice{
		if strings.Contains(text, key) {
			return true
		}
	}
	return false
}

func grab() []map[string]string{ 
	var wg sync.WaitGroup
	var jobs []map[string]string
	for i := 1; i < len(companies); i++ {
		wg.Add(1)
		companyName := strings.ToLower(companies[i])		
		go func(name string, pos int) { 
			defer wg.Done()
			x, err := goquery.NewDocument("https://jobs.dou.ua/companies/"+name+"/vacancies/")
			if err == nil {
				x.Find(".lt .l-vacancy").Each(func(i int, li *goquery.Selection) {
					jobName := strings.TrimSpace(li.Find(".title a").Text())
					description := strings.TrimSpace(li.Find(".sh-info").Text())
					if(isKeysContains(KEYS, jobName) || isKeysContains(KEYS, description)){
						link, _ := li.Find(".title a").Attr("href")
						element := map[string]string{
							"position":strconv.Itoa(pos),
							"jobName":jobName,
							"companyName": name,
							"description":description,
							"link":link,
						}
						jobs = append(jobs, element)
					}
				})
			}
			fmt.Print(name+" ")
		}(companyName, i)
		time.Sleep(100 * time.Millisecond)
	}
	//fmt.Println("Запущено потоков: ", i)
	wg.Wait()
	fmt.Print("\n")
	return jobs
}

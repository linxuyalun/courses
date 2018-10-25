package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var url = [9]string{
	"../EP4/data/buyer.0.0",
	"../EP4/data/buyer.1.1",
	"../EP4/data/good.0.0",
	"../EP4/data/good.1.1",
	"../EP4/data/good.2.2",
	"../EP4/data/order.0.0",
	"../EP4/data/order.0.3",
	"../EP4/data/order.1.1",
	"../EP4/data/order.2.2",
}

type itemList struct {
	List [][]string
}

// storeItem stroes every item of list from file into memory
func (l itemList) storeItem(i int) itemList {
	// read file
	file, err := os.Open(url[i])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// store item
	scanner := bufio.NewScanner(file)
	// In the given files of good and order, there is some bad smell.
	// transfer " " into "\t"
	if i == 0 || i == 1 {
		for scanner.Scan() {
			l.List = append(l.List, strings.Split(scanner.Text(), "\t"))
		}
	} else {
		r := strings.NewReplacer("    ", "\t", "   ", "\t", "  ", "\t", " ", "\t")
		for scanner.Scan() {
			text := r.Replace(scanner.Text())
			l.List = append(l.List, strings.Split(text, "\t"))
		}
	}
	return l
}

type order struct {
	orderID    string
	createTime int
	actualTime time.Time
	buyerID    string
	goodID     string
	remark     string
	amount     int
}

func normalizeOrder(storeOrder itemList) []order {
	var list []order
	for i := 0; i < len(storeOrder.List); i++ {
		var a order
		var err error
		a.orderID = strings.TrimPrefix(storeOrder.List[i][0], "orderid:")
		a.createTime, err = strconv.Atoi(strings.TrimPrefix(storeOrder.List[i][1], "createtime:"))
		if err != nil {
			panic(err)
		}
		// Calculate actual time
		a.actualTime = time.Unix(int64(a.createTime), 0)
		a.buyerID = strings.TrimPrefix(storeOrder.List[i][2], "buyerid:")
		a.goodID = strings.TrimPrefix(storeOrder.List[i][3], "goodid:")
		// Note that "remark" is not required
		if strings.HasPrefix(storeOrder.List[i][4], "remark:") {
			a.remark = strings.TrimPrefix(storeOrder.List[i][4], "remark:")
			a.amount, err = strconv.Atoi(strings.TrimPrefix(storeOrder.List[i][5], "amount:"))
			if err != nil {
				panic(err)
			}
		} else {
			a.amount, err = strconv.Atoi(strings.TrimPrefix(storeOrder.List[i][4], "amount:"))
			if err != nil {
				panic(err)
			}
		}
		list = append(list, a)
	}
	return list
}

type good struct {
	goodID   string
	salerID  string
	goodName string
	price    float64
}

func normalizeGood(storeGood itemList) []good {
	var list []good
	for i := 0; i < len(storeGood.List); i++ {
		var a good
		var err error
		a.goodID = strings.TrimPrefix(storeGood.List[i][0], "goodid:")
		a.salerID = strings.TrimPrefix(storeGood.List[i][1], "salerid:")
		a.goodName = strings.TrimPrefix(storeGood.List[i][2], "good_name:")
		// Note that "description" is not required.
		if strings.HasPrefix(storeGood.List[i][3], "price:") {
			a.price, err = strconv.ParseFloat(strings.TrimPrefix(storeGood.List[i][3], "price:"), 64)
			if err != nil {
				panic(err)
			}
		} else {
			a.price, err = strconv.ParseFloat(strings.TrimPrefix(storeGood.List[i][4], "price:"), 64)
			if err != nil {
				panic(err)
			}
		}
		list = append(list, a)
	}
	return list
}

type buyer struct {
	buyerID      string
	contactPhone string
	address      string
	buyerName    string
}

func normalizeBuyer(storeBuyer itemList) []buyer {
	var list []buyer
	for i := 0; i < len(storeBuyer.List); i++ {
		var a buyer
		a.buyerID = strings.TrimPrefix(storeBuyer.List[i][0], "buyerid:")
		// Note that both contactPhone and address are not required.
		for j := 1; j < 4; j++ {
			if strings.HasPrefix(storeBuyer.List[i][j], "contactphone:") {
				a.contactPhone = strings.TrimPrefix(storeBuyer.List[i][j], "contactphone:")
			} else if strings.HasPrefix(storeBuyer.List[i][j], "address:") {
				a.address = strings.TrimPrefix(storeBuyer.List[i][j], "address:")
			} else if strings.HasPrefix(storeBuyer.List[i][j], "buyername:") {
				a.buyerName = strings.TrimPrefix(storeBuyer.List[i][j], "buyername:")
				break
			}
		}
		list = append(list, a)
	}
	return list
}

func main() {

	var storeBuyer itemList
	var storeGood itemList
	var storeOrder itemList

	t := time.Now()
	// store item into list
	for i := 0; i < 2; i++ {
		storeBuyer = storeBuyer.storeItem(i)
	}
	for i := 2; i < 5; i++ {
		storeGood = storeGood.storeItem(i)
	}
	for i := 5; i < 9; i++ {
		storeOrder = storeOrder.storeItem(i)
	}

	orderlist := normalizeOrder(storeOrder)
	goodlist := normalizeGood(storeGood)
	buyerlist := normalizeBuyer(storeBuyer)

	var buyerConsumption = map[string]float64{}
	for _, v := range buyerlist {
		buyerConsumption[v.buyerID] = 0
	}

	var goodSales = map[string]int{}
	for _, v := range goodlist {
		goodSales[v.goodID] = 0
	}

	// Traversing orderlist to record buyerConsumption and goodSales
	for _, v := range orderlist {
		goodSales[v.goodID] += v.amount
		for _, v2 := range goodlist {
			if v2.goodID == v.goodID {
				buyerConsumption[v.goodID] += float64(v.amount) * v2.price
				break
			}
		}
	}

	// Rank top 3 buyer
	var firstBuyer string
	var secondBuyer string
	var thirdBuyer string
	for key, value := range buyerConsumption {
		if value > buyerConsumption[firstBuyer] {
			thirdBuyer = secondBuyer
			secondBuyer = firstBuyer
			firstBuyer = key
		} else if value > buyerConsumption[secondBuyer] {
			thirdBuyer = secondBuyer
			secondBuyer = key
		} else if value > buyerConsumption[thirdBuyer] {
			thirdBuyer = key
		}
	}
	fmt.Printf("Top 3 Buyers are: \n")
	fmt.Printf("%v : %v\n", firstBuyer, buyerConsumption[firstBuyer])
	fmt.Printf("%v : %v\n", secondBuyer, buyerConsumption[secondBuyer])
	fmt.Printf("%v : %v\n", thirdBuyer, buyerConsumption[thirdBuyer])

	// Rank top 3 good
	var firstGood string
	var secondGood string
	var thirdGood string
	for key, value := range goodSales {
		if value > goodSales[firstGood] {
			thirdGood = secondGood
			secondGood = firstGood
			firstGood = key
		} else if value > goodSales[secondGood] {
			thirdGood = secondGood
			secondGood = key
		} else if value > goodSales[thirdGood] {
			thirdGood = key
		}
	}
	fmt.Printf("Top 3 Goods are: \n")
	fmt.Printf("%v : %v\n", firstGood, goodSales[firstGood])
	fmt.Printf("%v : %v\n", secondGood, goodSales[secondGood])
	fmt.Printf("%v : %v\n", thirdGood, goodSales[thirdGood])

	fmt.Println(time.Now().Sub(t))

	type classifier struct {
		weenkends int
		weekdays  int
		day       int
		night     int
		buyerType string
	}
	var buyerConsumptionBehaviour = map[string]classifier{}
	for _, v := range buyerlist {
		var init classifier
		buyerConsumptionBehaviour[v.buyerID] = init
	}

	for _, v := range orderlist {
		a := buyerConsumptionBehaviour[v.buyerID]
		// Order numbers on weekends and weekdays
		if v.actualTime.Weekday().String() == "Saturday" || v.actualTime.Weekday().String() == "Sunday" {
			a.weenkends++
			buyerConsumptionBehaviour[v.buyerID] = a
		} else {
			a.weekdays++
			buyerConsumptionBehaviour[v.buyerID] = a
		}
		// Order numbers on day and night
		if v.actualTime.Hour() > 18 || v.actualTime.Hour() < 9 {
			a.night++
			buyerConsumptionBehaviour[v.buyerID] = a
		} else {
			a.day++
			buyerConsumptionBehaviour[v.buyerID] = a
		}
	}

	// Record the type of buyer by order numbers in classifier
	// A represents (weekends/2 > weekdays/5) && (night > day)
	// B represents (weekends/2 > weekdays/5) && (night <= day)
	// C represents (weekends/2 <= weekdays/5) && (night > day)
	// D represents (weekends/2 <= weekdays/5) && (night <= day)
	for key, value := range buyerConsumptionBehaviour {
		temp := value
		if float64(value.weenkends)/2 > float64(value.weekdays)/5 {
			if value.night > value.day {
				temp.buyerType = "A"
			} else {
				temp.buyerType = "B"
			}
		} else {
			if value.night > value.day {
				temp.buyerType = "C"
			} else {
				temp.buyerType = "D"
			}
		}
		buyerConsumptionBehaviour[key] = temp
	}

	// Give a buyer id, judge the type of buyer and find the buyers whose consumption behaviour are the same
	buyerID := "ap-bf11-8ff973e02aaf"
	fmt.Printf("The buyer %v is type %v \n", buyerID, buyerConsumptionBehaviour[buyerID].buyerType)
	fmt.Printf("The fllowing buyers have the same consumption behaviour of %v :\n", buyerID)
	for key, value := range buyerConsumptionBehaviour {
		if value == buyerConsumptionBehaviour[buyerID] {
			if key != buyerID {
				fmt.Println(key)
			}
		}
	}

}

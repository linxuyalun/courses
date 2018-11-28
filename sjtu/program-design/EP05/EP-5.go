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
	goodID      string
	salerID     string
	goodName    string
	description string
	price       float64
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
			a.description = strings.TrimPrefix(storeGood.List[i][3], "description:")
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
	// (may use in future)buyerlist := normalizeBuyer(storeBuyer)

	var consumption = map[string]float64{
		"Monday":    0,
		"Tuesday":   0,
		"Wednesday": 0,
		"Thursday":  0,
		"Friday":    0,
		"Saturday":  0,
		"Sunday":    0,
	}

	var orderNum = map[string]int{
		"Monday":    0,
		"Tuesday":   0,
		"Wednesday": 0,
		"Thursday":  0,
		"Friday":    0,
		"Saturday":  0,
		"Sunday":    0,
	}

	buyerID := "wx-805a-89cb83fd5551"
	for i := 0; i < len(orderlist); i++ {
		if buyerID == orderlist[i].buyerID {
			// get price
			var price float64
			for j := 0; j < len(goodlist); j++ {
				if orderlist[i].goodID == goodlist[j].goodID {
					price = goodlist[j].price
					break
				}
			}

			sum := price * float64(orderlist[i].amount)

			// get weekday
			t := time.Unix(int64(orderlist[i].createTime), 0)
			weekday := t.Weekday().String()
			// Record data
			consumption[weekday] += sum
			orderNum[weekday]++
		}
	}

	for key, value := range consumption {
		fmt.Printf("%v:\nOrdernumber: %v, Sum: %v \n", key, orderNum[key], value)
	}
}

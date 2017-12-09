package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("enter a number: ")
	for scanner.Scan() {
		line := scanner.Text()
		totOne, totTwo, err := parseLine(line)
		if err != nil {
			fmt.Printf("bad line: %s\n", err)
			fmt.Print("enter a number: ")
			continue
		}

		fmt.Printf("scores for %s are %d and %d\n", line, totOne, totTwo)
		fmt.Print("enter a number: ")
	}

	if scanner.Err() != nil {
		os.Stderr.WriteString(fmt.Sprintf("scan error %s", scanner.Err))
	}
}

func parseLine(line string) (totOne, totTwo uint64, err error) {

	nums, err := strToArray(line)
	if err != nil {
		return 0, 0, err
	}

	var numCount = len(nums)
	if numCount%2 != 0 {
		return 0, 0, errors.New("odd number of numbers")
	}

	var offsetOne = 1
	var offsetTwo = numCount / 2
	for i := 0; i < numCount; i++ {
		partnerOne := nums[(i+offsetOne)%numCount]
		partnerTwo := nums[(i+offsetTwo)%numCount]
		//fmt.Printf("%d has partner at %d\n", i, partnerIdx)
		if nums[i] == partnerOne {
			totOne += uint64(nums[i])
		}
		if nums[i] == partnerTwo {
			totTwo += uint64(nums[i])
		}
	}

	return totOne, totTwo, nil
}

func strToArray(line string) ([]uint8, error) {
	digits := strings.Split(line, "")
	nums := make([]uint8, len(digits))

	for idx, digit := range digits {
		num, err := strconv.ParseUint(digit, 10, 8)
		if err != nil {
			return nil, errors.New("not a number")
		}
		nums[idx] = uint8(num)
	}
	return nums, nil
}

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
		tot, err := parseLine(line)
		if err != nil {
			fmt.Printf("bad line: %s\n", err)
			fmt.Print("enter a number: ")
			continue
		}

		fmt.Printf("score for %s is %d\n", line, tot)
		fmt.Print("enter a number: ")
	}

	if scanner.Err() != nil {
		os.Stderr.WriteString(fmt.Sprintf("scan error %s", scanner.Err))
	}
}

func parseLine(line string) (uint64, error) {

	nums, err := strToArray(line)
	if err != nil {
		return 0, err
	}

	var tot uint64
	var prevNum uint8
	for _, num := range nums {
		if num == prevNum {
			tot += uint64(num)
		}
		prevNum = num
	}

	if nums[0] == nums[len(nums)-1] {
		// wrap match
		tot += uint64(nums[0])
	}

	return tot, nil
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

package main

import "fmt"


// 136. Single Number
func singleNumber(nums []int) int {
    var countMap = make(map[int]int)
    for i := 0; i < len(nums);i++ {
        countMap[nums[i]]++
    }
    
    var ans int
    for i, count := range countMap {
        if count == 1 {
            ans = i
            break
        }
    }

    return ans
}


// 9. Palindrome Number
func isPalindrome(x int) bool {
    if x < 0 {
        return false
    }

    var original = x
    var reversed = 0

    for x > 0 {
        reversed = reversed * 10 + x % 10
        x = x / 10
    }

	return original == reversed
}

// 20. Valid Parentheses
func isValid(s string) bool {
    var stack = make([]rune, 0)

    for _, char := range s {
        if char == '(' || char == '[' || char == '{' {
            stack = append(stack, char)
        } else {
            if len(stack) == 0 {
                return false
            }	

			top := stack[len(stack)-1]

			if char == ')' && top == '(' || char == ']' && top == '[' || char == '}' && top == '{' {
				stack = stack[:len(stack)-1]
			} else {
				return false
			}
        }
	}

	return len(stack) == 0;
}

// 14. Longest Common Prefix
func longestCommonPrefix(strs []string) string {
    if len(strs) == 0 {
        return ""
    }

    prefix := strs[0]

    for i := 1; i < len(strs); i++ {
        curr := strs[i]
        j := 0
        for j < len(prefix) && j < len(curr) && prefix[j] == curr[j] {
            j++
        }
        prefix = prefix[:j]
        if prefix == "" {
            return ""
        }
    }

    return prefix
}


// 66. Plus One
func plusOne(digits []int) []int {
    for i := len(digits) - 1; i >= 0; i-- {
        if digits[i] < 9 {
            digits[i]++
            return digits
        }
        digits[i] = 0
    }
    return append([]int{1}, digits...)
}

// 26. Remove Duplicates from Sorted Array
func removeDuplicates(nums []int) int {
    if len(nums) == 0 {
        return 0
    }
    var i = 0
    for j := 1; j < len(nums); j++ {
        if nums[j] != nums[i] {
            i++
            nums[i] = nums[j]
        }
    }
    return i + 1
}


// 56. Merge Intervals
func merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return [][]int{}
	}

	
	for i := 0; i < len(intervals)-1; i++ {
		minIndex := i
		for j := i + 1; j < len(intervals); j++ {
			if intervals[j][0] < intervals[minIndex][0] {
				minIndex = j
			}
		}
		
		if minIndex != i {
			temp := intervals[i]
			intervals[i] = intervals[minIndex]
			intervals[minIndex] = temp
		}
	}

	result := make([][]int, 0)
	current := make([]int, 2)
	current[0] = intervals[0][0]
	current[1] = intervals[0][1]

	for i := 1; i < len(intervals); i++ {
		start := intervals[i][0]
		end := intervals[i][1]

		
		if start <= current[1] {
		
			if end > current[1] {
				current[1] = end
			}
		} else {
	
			result = append(result, []int{current[0], current[1]})
			current[0] = start
			current[1] = end
		}
	}


	result = append(result, []int{current[0], current[1]})

	return result
}



func twoSum(nums []int, target int) []int {
    var result = make([]int, 0)

    for i := 0; i < len(nums); i++ {
        for j := i + 1; j < len(nums); j++ {
            if nums[i] + nums[j] == target {
                result = append(result, i, j)
                return result
            }
        }
    }
	return result
}



func main() {
		// nums := []int{2,2,3,2}
		// fmt.Println(singleNumber(nums))
		//fmt.Println(isPalindrome(121))
		//fmt.Println(isValid("()"))
		//fmt.Println(longestCommonPrefix([]string{"flower","flow","flight"}))
		//fmt.Println(plusOne([]int{9,9,9}))
		//fmt.Println(removeDuplicates([]int{1,1,2}))
		//fmt.Println(merge([][]int{{1,3},{2,6},{8,10},{15,18}}))
		fmt.Println(twoSum([]int{2,7,11,15}, 9))


}



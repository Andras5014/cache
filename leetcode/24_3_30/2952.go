package _4_3_30

import (
	"sort"
)

func main() {

}
func minimumAddedCoins(coins []int, target int) (ans int) {
	sort.Ints(coins)
	for i, s := 0, 1; s <= target; i++ {
		if i < len(coins) && coins[i] <= s {
			s += coins[i]
		} else {
			ans++
			s *= 2
		}
	}
	return ans
}

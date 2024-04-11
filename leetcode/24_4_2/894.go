package main



 type TreeNode struct {
     Val int
   Left *TreeNode
     Right *TreeNode
 }

func allPossibleFBT(n int) []*TreeNode {
	trees := make([]*TreeNode, 0)
	if n%2==0 {
	return 	trees
	}
	if n==1 {
		return append(trees,&TreeNode{Val:0})
	}
	for i:=1;i<n;i+=2 {
		left := allPossibleFBT(i)
		right := allPossibleFBT(n-1-i)
		for _,l:=range left {
			for _,r:=range right {
				trees = append(trees,&TreeNode{Val:0,Left:l,Right:r})
			}
		}
	}
	return trees
}
func main() {

}

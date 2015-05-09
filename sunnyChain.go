/**
 这是一个关于舞者搭配的题目.需要使用haskell语言实现,折腾了半天,
 现在使用go语言比较简单地实现了相关功能
 作者:徐洋
 时间:2015年 5月8日
 版权 MIT协议
**/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

var leaderNum int
var followerNum int
var dance [8]string

type inviteCard struct {
	leader     int
	follower   int
	danceQuene int
}

func LeaderProc(ch chan inviteCard, id int, followerCh []chan inviteCard, mainCh []chan int) {
	buf := [8]int{-1, -1, -1, -1, -1, -1, -1, -1}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for step := 0; step < 8; step++ {

		start := r.Intn(followerNum)
		for i := 0; i < followerNum; i++ {
			inviteFollower := (start + i) % followerNum
			myInviteCard := new(inviteCard)
			myInviteCard.danceQuene = step
			myInviteCard.follower = inviteFollower
			myInviteCard.leader = id
			followerCh[inviteFollower] <- *myInviteCard
			time.Sleep(1e6)
			myReadInviteCard := <-ch
			//接受
			if myReadInviteCard.leader != -1 {
				//更新
				buf[step] = myReadInviteCard.follower
				break
			}

		}
	}

	mainCh[id] <- id
	<-mainCh[id]

	fmt.Println("Leader:", id)

	for i := 0; i < 8; i++ {

		fmt.Println(dance[i], "\t with", buf[i])
	}

}

func canIDanceWithLeader(leader int, step int, buf [8]int) (isOK bool) {
	if buf[step] == -1 {
		//计算已经和leader跳过几次了
		count := 0
		for i := 0; i < 8; i++ {

			if buf[i] == leader {
				count++
			}
		}
		if count < 2 {
			return true
		}

	}
	return false
}

func FollowerProc(ch chan inviteCard, id int, leaderCh []chan inviteCard) {
	buf := [8]int{-1, -1, -1, -1, -1, -1, -1, -1} //个人资料库
	for {
		myReadInviteCard := <-ch //收到邀请
		theLead := myReadInviteCard.leader
		theStep := myReadInviteCard.danceQuene
		theFollower := myReadInviteCard.follower

		if id != theFollower {
			fmt.Println("erro!!!!")
		}

		if canIDanceWithLeader(theLead, theStep, buf) {
			buf[theStep] = theLead
			leaderCh[theLead] <- myReadInviteCard
		} else {
			myReadInviteCard.leader = -1
			leaderCh[theLead] <- myReadInviteCard
		}
	}

}

//产生 全局的变量
func Myinit() {
	leaderNum = 8
	followerNum = 15
	dance = [8]string{"wal", "rumba", "samba", "cha cha", "sunny", "hello", "world", "I AM OK"}
	return
}
func main() {

	Myinit()
	//每个 leader 和follower 的chan
	leaderChan := make([]chan inviteCard, leaderNum)
	followerChan := make([]chan inviteCard, followerNum)
	mainChan := make([]chan int, leaderNum)
	//创建 leader 的线程
	for i := 0; i < leaderNum; i++ {
		leaderChan[i] = make(chan inviteCard, 100)
		mainChan[i] = make(chan int)
	}

	//创建 follower 个 线程
	for i := 0; i < followerNum; i++ {

		followerChan[i] = make(chan inviteCard, 100)

	}

	//创建 follower 线程
	for i := 0; i < followerNum; i++ {

		go FollowerProc(followerChan[i], i, leaderChan)
	}
	//创建 leader 线程
	for i := 0; i < leaderNum; i++ {

		go LeaderProc(leaderChan[i], i, followerChan, mainChan)
	}
	//等待 结束
	nCount := 0

	for nCount < leaderNum {
		<-mainChan[nCount]
		time.Sleep(1e7)
		mainChan[nCount] <- -1
		nCount++

	}

	time.Sleep(1e7)
	//结束打印
	fmt.Println("finish")

}

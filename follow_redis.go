package main

import "feedscenter/models/rdsmodel"

// UEM_user_followers_v0:4329377 用户的粉丝列表
// UEM_user_followings_v0:4329299

func test() {
    c := rdsmodel.Connection()
    defer c.Close()

    // uem redis db: 11
    c.Do("SELECT", 11)
    c.Do("zRange", "UEM_user_followers_v0:4329377", 0, -1)
}

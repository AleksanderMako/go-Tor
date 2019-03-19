"use strict";
var redis = require("redis");
var Err = require("../../types/error/error.js")

function createRedisClient () {
    const redisUrl = "redis://redis-persistance:6379";
    const client= redis.createClient(redisUrl);
    return client;
}

function put (redisClient,data) {
// standardize all data to objects 
    redisClient.on("error", function (err) {
        console.log("Error in redis service put operation  " + err);
    });
    const dataObject = {
        redisData: data 
    }
    redisClient.set("testKey", JSON.stringify(dataObject));
  
   // redisClient.quit();
}

function get(redisClient,key) {

    return new Promise((resolve,reject)=>{

         redisClient.get(key, (err, data) => {

            if (err != null) {

            //    const redisErr = new Err(true, err)
            console.log("Redis service err :", err);
                reject (err);
            }
            const dataObject = JSON.parse(data);
             resolve(dataObject.redisData);

        });

    })

}

module.exports = {
    createRedisClient,
    put,
    get
}
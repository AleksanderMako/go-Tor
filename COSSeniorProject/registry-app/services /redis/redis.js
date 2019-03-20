"use strict";
var redis = require("redis");
var Err = require("../../types/error/error.js")

function createRedisClient() {
    const redisUrl = "redis://redis-persistance:6379";
    const client = redis.createClient(redisUrl);
    return client;
}

function put(redisClient, data, key) {
    // standardize all data to objects 
    redisClient.on("error", function (err) {
        console.log("Error in redis service put operation  " + err);
    });
    const dataObject = {
        redisData: data
    }
    redisClient.set(key, JSON.stringify(dataObject));

    redisClient.quit();
}

function get(redisClient, key) {

    return new Promise((resolve, reject) => {

        redisClient.get(key, (err, data) => {

            if (err != null) {

                //    const redisErr = new Err(true, err)
                console.log("Redis service err :", err);
                reject(err);
            }
            const dataObject = JSON.parse(data);
            resolve(dataObject.redisData);

        });

    })

}

function putInSet(redisClient, setKey, dataKey, data) {

    redisClient.on("error", function (err) {
        console.log("Error in redis service put operation  " + err);
    });
    const dataObject = {
        redisData: data
    }
    redisClient.hset(setKey, dataKey, JSON.stringify(dataObject))
    redisClient.quit();


}

function getSetKeys(redisClient, setKey) {

    return new Promise((resolve, reject) => {

        redisClient.hkeys(setKey, (err, replies) => {

            if (err != null || !replies) {
                reject(err)
            }
            resolve(replies)
        });
    });
}

function getItemInSet(redisClient, setKey, dataKey) {

    return new Promise((resolve, reject) => {

        redisClient.hget(setKey, dataKey, (err, data) => {

            if (err != null || !data) {
                reject(err)
            }
            const dataObj = JSON.parse(data);
            resolve(dataObj)
        });
    })
}

async function getSetElements(redisClient, setKey) {

    var setItemsArray = new Array();
    const setKeys = await getSetKeys(redisClient, setKey)

    setKeys.forEach(async (key) => {

        let setItem = await getItemInSet(redisClient, setKey, key)
        setItemsArray.push(setItem)
    });
    return setItemsArray;
}
module.exports = {
    createRedisClient,
    put,
    get,
    getItemInSet,
    getSetKeys,
    putInSet,
    getSetElements
}
var express = require('express');
var router = express.Router();
var Response = require("../types/response/response");
var redisService = require("../services /redis/redis");

router.post("/", (req, res) => {
    // peers register under the peer set 
    const serviceAddress = req.body.id;
    const ips = req.body.ips;
    const keyWords = req.body.keyWords;
    console.log("body of request is ", req.body);
    if (!serviceAddress) {
        return res.status(500).send("Missing service id ")
    }
    if (!ips) {
        return res.status(500).send("Missing service introduction points  ")
    }
    if (!keyWords) {
        return res.status(500).send("Missing keyWords")
    }

    const redisClient = redisService.createRedisClient();
    const setKey = "descriptor"
    redisService.putInSet(redisClient, setKey, serviceAddress, ips);
    ips.forEach(ip => {
        const redisClient = redisService.createRedisClient();
        redisService.putInSet(redisClient, "peers", ip, "occupied");

    });
    for (let i = 0; i < keyWords.length; i++) {

        const redisClient = redisService.createRedisClient();
        redisService.putInSet(redisClient, keyWords[i], serviceAddress, ips);

    }
    return res.send("successfully added service description")

});
router.get("/:serviceID", async (req, res) => {

    const serviceID = req.params.serviceID;
    if (!serviceID) {
        return res.status(500).send("MISSING SERVICE ID ")
    }
    const redisClient = redisService.createRedisClient();
    const setKey = "descriptor";
    const redisDataObj = await redisService.getItemInSet(redisClient, setKey, serviceID);
    return res.json(redisDataObj);


});
router.get("/:keyWord", async (req, res) => {

    const keyWord = req.params.keyWord;
    if (!keyWord) {
        return res.status(500).send("Missing keyWord")
    }
    // console.log("keyword", keyWord);
    // let redisClient = redisService.createRedisClient();
    // const setItems = await redisService.getSetKeys(redisClient, keyWord);
    // // const servicDescriptors = await redisService.getServiceDescriptor(setItems, keyWord);
    // const response = {
    //     serviceDescriptors: setItems,
    // }
    return res.json(keyWord);
});

module.exports = router;
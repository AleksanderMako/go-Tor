var redis = require("../../services /redis/redis");
var chai = require("chai");
var mocha = require("mocha");
describe("redis service initial tests !", () => {

    it("should put data correctly ", async () => {

        // Arrange 
        const test = {
            redisData: "this_data"
        }
        let client = redis.createRedisClient();
        const key = "k1"
        //Act 
        redis.put(client, test, key);

        //Assert 
        client = redis.createRedisClient();
        const data = await redis.get(client, key);
        chai.expect(data.redisData).equal(test.redisData);

    });

    it("should put item in set  correctly", async () => {

        //Arrange 
        let client = redis.createRedisClient()
        const setKey = "set1"
        const dataKey = "data1"
        const test = "this_data"


        //Act 
        redis.putInSet(client, setKey, dataKey, test)

        //Assert 
        client = redis.createRedisClient()
        const data = await redis.getItemInSet(client, setKey, dataKey)
        chai.expect(data.redisData).equal(test);

    });

    // TODO:refactor to quit the client via some method called at the end of the test
    it("should get set correctly", async () => {

        //Arrange
        let client = redis.createRedisClient()
        const setKey = "set2"
        const dataKey1 = "data1"
        const dataKey2 = "data2"
        const dataKey3 = "data3"
        const td1 = "testData1"
        const td2 = "testData2"
        const td3 = "testData3"


        redis.putInSet(client, setKey, dataKey1, td1);

        client = redis.createRedisClient()
        redis.putInSet(client, setKey, dataKey2, td2);

        client = redis.createRedisClient()
        redis.putInSet(client, setKey, dataKey3, td3);

        //Act 
        client = redis.createRedisClient()

        const set = await redis.getSetKeys(client, setKey)

        chai.expect(set).to.deep.equal([dataKey1, dataKey2, dataKey3]);

    });


})
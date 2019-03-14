"use strict";

module.exports = class Response {

    constructor() {

    }

    getOpStatus (opSts, hassErr, message ) {
        return new Response={
            operationStatus: opSts,
            hassErr : hassErr,
            message: message,
        }
    }
}
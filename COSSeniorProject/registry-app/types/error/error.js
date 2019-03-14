"use strict";

module.exports = class ERR {

    constructor ( hasErr, message ) {
        this.hasErr = true;
        this.message = message;
    }
}
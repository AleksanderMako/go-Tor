"use strict";

function getOpStatus(opSts, hassErr, message) {

    var response = {
        operationStatut:opSts,
        hassErr:hassErr,
        message:message,
    }
    return response;
}
module.exports = {
    getOpStatus
}
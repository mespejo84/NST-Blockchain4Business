'use strict';
const shim = require('fabric-shim');
const util = require('util');

class Chaincode {

    get chaincodeFunctions() {
        return {
            'create': this.createLicenses,
            'transfer': this.transferLicense,
            'query': this.queryLicense
        }
    };

    async Init(stum) {
        console.info("****Init Chaincode****");
        return shim.success();
    }

    async Invoke(stub) {
        console.info("****Invoke Chaincode****");

        let ret = stub.getFunctionAndParameters();
        console.info("ret", ret);

        let method = this.chaincodeFunctions[ret.fcn];
        if (!method) {
            console.error("Invalid function name");
            throw new Error("Invalid function name " + ret.fcn);
        }
        try {
            let payload = await method(stub, ret.params);
            return shim.success(payload);
        } catch (err) {
            console.error(err);
            return shim.error(err);
        }
    }

    async createLicenses(stub, args) {
        console.info("******Create licenses BEGIN******");

        let licenses = [
            {licenseId: "abc123", Name: "MyWordProcessor", organization: "MicroSuave", licensedTo: "Billy"},
            {licenseId: "abc456", name: "MyWordProcessor", organization: "MicroSuave", licensedTo: "Billy"},
            {licenseId: "abc789", name: "MyWordProcessor", organization: "MicroSuave", licensedTo: "Billy"},
            {licenseId: "abc987", name: "MyWordProcessor", organization: "MicroSuave", licensedTo: "Billy"},
            {licenseId: "abc654", name: "MyWordProcessor", organization: "MicroSuave", licensedTo: "Billy"},
            {licenseId: "abc321", name: "MyWordProcessor", organization: "MicroSuave", licensedTo: "Billy"}
        ];

        for (let license of licenses) {
            await stub.putState(license.licenseId, Buffer.from(JSON.stringify(license)));
            console.info("License Added: ", license.licenseId);
        }

        console.info("******Create licenses END******");
    }

    async transferLicense(stub, args) {
        console.info("******Transfer license BEGIN******");
        if (args.length != 2) {
            throw new Error("Incorrect number of arguments");
        }

        let [ licenseId, transferTo ] = args;

        let licenseAsBytes = await stub.getState(licenseId);

        if (!licenseAsBytes || licenseAsBytes.toString().length <= 0) {
            throw new Error("ailed to get current state for license ", licenseId);
        }

        let license = JSON.parse(licenseAsBytes);
        if (license.licensedTo !== "Billy") {
            throw new Error("The license has been transferred to ", license.licensedTo);
        }

        license.licensedTo = transferTo;

        await stub.putState(licenseId, Buffer.from(JSON.stringify(license)));

        console.info("******Transfer license END******");
    }

    async queryLicense(stub, args) {
        console.info("******Query license BEGIN******");
        if (args.length != 1) {
            throw new Error("Incorrect number of arguments");
        }

        let licenseId = args[0];

        let licenseAsBytes = await stub.getState(licenseId);
        if (!licenseAsBytes || licenseAsBytes.toString().length <= 0) {
            throw new Error("ailed to get current state for license ", licenseId);
        }
        
        let license = JSON.parse(licenseAsBytes);
        console.info("License Data: ", license);

        console.info("******Query license END******");

        return licenseAsBytes;
    }
};

shim.start(new Chaincode());
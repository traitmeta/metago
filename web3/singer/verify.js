//创建web3对象
var Web3 = require('web3');
var sigUtil = require("eth-sig-util")
var provider = new Web3.providers.HttpProvider("http://localhost:7545");
var web3 = new Web3(provider);

var json = require("../build/contracts/Demo.json");
var contractAddr = '';

var account = "";
var account_to = "";
var privateKey = "";
var privateKeyHex = Buffer.from(privateKey, 'hex')

var demoContract = new web3.eth.Contract(json['abi'], contractAddr);

//获取链ID
demoContract.methods.getChainId().call({ from: account }, function (error, result) {
    if (error) {
        console.log(error);
    }
    console.log("getChainId:", result);
});

//V4签名
const typedData = {
    types: {
        EIP712Domain: [
            { name: 'name', type: 'string' },
            { name: 'version', type: 'string' },
            { name: 'chainId', type: 'uint256' },
            { name: 'verifyingContract', type: 'address' },
        ],
        Mail: [
            { name: 'from', type: 'address' },
            { name: 'to', type: 'address' },
            { name: 'value', type: 'uint256' },
        ],
    },
    domain: {
        name: 'Demo',
        version: '1.0',
        chainId: 1,
        verifyingContract: contractAddr,
    },
    primaryType: 'Mail',
    message: {
        from: account,
        to: account_to,
        value: 12345,
    },
}

//V4签名
var signature = sigUtil.signTypedData_v4(privateKeyHex, { data: typedData })
console.log("signature: ", signature)

//V4验签
const recovered = sigUtil.recoverTypedSignature_v4({
    data: typedData,
    sig: signature,
});
console.log("recovered: ", recovered)

//合约V4验签
demoContract.methods.verify(typedData.message.from, typedData.message.to, typedData.message.value, signature).call({ from: account }, function (error, result) {
    if (error) {
        console.log(error);
    }
    console.log("verify: ", result);
});


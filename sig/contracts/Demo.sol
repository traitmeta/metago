//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.0;

import "hardhat/console.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/cryptography/draft-EIP712.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
 
contract Demo is EIP712 {
    constructor(string memory name, string memory version) EIP712(name, version) {
    }
 
    //获取签名人（V4）
    function recoverV4(
        address from,
        address to,
        uint256 value,
        bytes memory signature
    ) public view returns (address) {
        bytes32 digest = _hashTypedDataV4(keccak256(abi.encode(
        keccak256("Mail(address from,address to,uint256 value)"),
            from,
            to,
            value
        )));
        return ECDSA.recover(digest, signature);
    }
 
    //验签
    function verify(
        address from,
        address to,
        uint256 value,
        bytes memory signature
    ) public view returns (bool) {
        address signer = recoverV4(from, to, value, signature);
        return signer == from;
    }
 
    function getChainId() public view returns(uint256) {
        return block.chainid;
    }
}
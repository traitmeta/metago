// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

contract WETH {
    event Deposit(address indexed from, uint256 amount);
    event Withdraw(address indexed to, uint256 amount);

    event Transfer(address indexed _from, address indexed _to, uint256 _value);
    event Approval(
        address indexed _owner,
        address indexed _spender,
        uint256 _value
    );
}

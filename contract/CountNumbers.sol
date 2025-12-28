// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract CountNumbers {

   uint256 public  counter;

   event count(address send);


constructor(uint256 value){
    counter=value;
}

  function add() public returns (uint) {
    counter++;
    emit count(msg.sender);
    return counter;
  }

    function getCount() public view returns (uint256) {
        return counter;
  }




}
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;


contract TransferNFT {
  event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);
  event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value);
  event TransferBatch(
      address indexed operator,
      address indexed from,
      address indexed to,
      uint256[] ids,
      uint256[] values
  );
}
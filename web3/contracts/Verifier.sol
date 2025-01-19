// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "@openzeppelin/contracts/access/Ownable.sol";

contract Verifier is Ownable {
    bytes32 public merkleRoot;
    event MerkleRootUpdated(bytes32 oldRoot, bytes32 newRoot);

    constructor(bytes32 _merkleRoot) Ownable(msg.sender) {
        merkleRoot = _merkleRoot;
    }

    function verify(
        address user,
        bytes32[] calldata proof
    ) public view returns (bool) {
        // If there is only one leaf in the entire tree (i.e. proof is empty),
        // the leaf itself should match the root.
        if (proof.length == 0) {
            bytes32 leaf = keccak256(abi.encodePacked(user));
            return leaf == merkleRoot;
        }

        // Normal multi-leaf case
        bytes32 computedHash = keccak256(abi.encodePacked(user));

        for (uint256 i = 0; i < proof.length; i++) {
            bytes32 proofElement = proof[i];
            bytes32 tempHash;

            assembly {
                let memPtr := mload(0x40)

                switch lt(computedHash, proofElement)
                case 1 {
                    mstore(memPtr, computedHash)
                    mstore(add(memPtr, 32), proofElement)
                }
                default {
                    mstore(memPtr, proofElement)
                    mstore(add(memPtr, 32), computedHash)
                }

                tempHash := keccak256(memPtr, 64)
            }
            computedHash = tempHash;
        }

        return computedHash == merkleRoot;
    }

    function updateMerkleRoot(bytes32 _newRoot) public onlyOwner {
        bytes32 oldRoot = merkleRoot;
        merkleRoot = _newRoot;
        emit MerkleRootUpdated(oldRoot, _newRoot); // for tracking purposes
    }
}

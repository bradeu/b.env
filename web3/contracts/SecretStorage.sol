// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "./Verifier.sol";

contract SecretStorage {
    // struct
    struct ApiLots {
        address lotAddress;
        string encryptedApiKey;
    }

    Verifier public verifier;
    mapping(address => ApiLots) public lots;

    event EncryptedApiKeyStored(address indexed user, string encryptedApiKey);

    constructor(address _verifierAddress) {
        verifier = Verifier(_verifierAddress);
    }

    //function made public
    function storeEncryptedApiKey(
        address targetAddress,
        string memory _encryptedApiKey,
        bytes32[] calldata proof
    ) public {
        require(
            verifier.verify(msg.sender, proof),
            "Not authorized to store API keys"
        );
        // ensure API key hasn't been stored for given address
        require(
            bytes(lots[targetAddress].encryptedApiKey).length == 0,
            "API key already stored"
        );

        lots[targetAddress] = ApiLots({
            lotAddress: targetAddress,
            encryptedApiKey: _encryptedApiKey
        });

        emit EncryptedApiKeyStored(targetAddress, _encryptedApiKey);
    }

    // get is made public
    function getEncryptedApiKeyForAddress(
        address targetAddress,
        bytes32[] calldata proof
    ) public view returns (string memory) {
        require(
            verifier.verify(msg.sender, proof),
            "Not authorized to access this API key"
        );
        return lots[targetAddress].encryptedApiKey;
    }
}

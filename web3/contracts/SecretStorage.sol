// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

contract SecretStorage {
    // struct
    struct ApiLots {
        address lotAddress;
        string encryptedApiKey;
    }

    mapping(address => ApiLots) public lots;

    event EncryptedApiKeyStored(address indexed user, string encryptedApiKey);

    //function made public
    function storeEncryptedApiKey(
        address targetAddress,
        string memory _encryptedApiKey
    ) public {
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
        address targetAddress
    ) public view returns (string memory) {
        return lots[targetAddress].encryptedApiKey;
    }
}

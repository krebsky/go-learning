// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title AuctionNFT
 * @dev ERC721 NFT 合约，支持铸造和转移
 */
contract AuctionNFT is ERC721URIStorage, Ownable {
    uint256 private _tokenIds;
    
    // 存储每个 token ID 对应的元数据 URI
    mapping(uint256 => string) private _tokenURIs;
    
    constructor(string memory name, string memory symbol) ERC721(name, symbol) Ownable(msg.sender) {}
    
    /**
     * @dev 铸造新的 NFT
     * @param recipient NFT 接收者地址
     * @param tokenUri NFT 元数据 URI
     * @return 新铸造的 token ID
     */
    function mintNFT(address recipient, string memory tokenUri) public returns (uint256) {
        _tokenIds++;
        uint256 newTokenId = _tokenIds;
        
        _mint(recipient, newTokenId);
        _setTokenURI(newTokenId, tokenUri);
        
        return newTokenId;
    }
    
    /**
     * @dev 批量铸造 NFT
     * @param recipient NFT 接收者地址
     * @param tokenUris NFT 元数据 URI 数组
     * @return tokenIds 新铸造的 token ID 数组
     */
    function batchMintNFT(address recipient, string[] memory tokenUris) public returns (uint256[] memory) {
        uint256[] memory tokenIds = new uint256[](tokenUris.length);
        
        for (uint256 i = 0; i < tokenUris.length; i++) {
            tokenIds[i] = mintNFT(recipient, tokenUris[i]);
        }
        
        return tokenIds;
    }
    
    /**
     * @dev 获取当前总供应量
     * @return 已铸造的 NFT 总数
     */
    function totalSupply() public view returns (uint256) {
        return _tokenIds;
    }
    
    /**
     * @dev 获取下一个 token ID
     * @return 下一个 token ID
     */
    function nextTokenId() public view returns (uint256) {
        return _tokenIds + 1;
    }
}


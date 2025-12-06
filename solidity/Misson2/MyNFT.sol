// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";


contract MyNFT is ERC721 {
    using Counters for Counters.Counter;
    
    // 用于跟踪 token ID 的计数器
    Counters.Counter private _tokenIds;
    
    // 存储每个 token ID 对应的 URI（元数据链接）
    mapping(uint256 => string) private _tokenURIs;
    
    

    constructor(string memory name, string memory symbol) ERC721(name, symbol) {
    }
    

    function mintNFT(address recipient, string memory tokenUri) public  returns (uint256) {
        _tokenIds.increment();
        
        uint256 newTokenId = _tokenIds.current();
        
        _mint(recipient, newTokenId);
        
        _tokenURIs[newTokenId] = tokenUri;
        
        
        return newTokenId;
    }
    

    function tokenURI(uint256 tokenId) public  view  override virtual returns (string memory) 
    {
        require(_exists(tokenId), "Error: URI query for nonexistent token");
        return _tokenURIs[tokenId];
    }
    

    function totalSupply() public view returns (uint256) {
        return _tokenIds.current();
    }
    

    function nextTokenId() public view returns (uint256) {
        return _tokenIds.current() + 1;
    }

    function _exists(uint tokenId) private view  returns(bool){
        return _ownerOf(tokenId) != address(0);
    }
}
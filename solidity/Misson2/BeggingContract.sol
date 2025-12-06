// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;



contract BeggingContract  {

    mapping(address => uint256) private _donations;
    
    address[] private _donors;
    
    address private _owner;

    uint256 public totalDonations;
    
    uint256 public donationStartTime;
    
    uint256 public donationEndTime;
    
    bool public timeLimitEnabled;
    
    // 事件：记录每次捐赠
    event Donation(address indexed donor, uint256 amount, uint256 timestamp);
    

    modifier onlyOwner(){
        require(msg.sender == _owner, "not owner");
        _;
    }
    
    constructor(bool _timeLimitEnabled, uint256 _startTime, uint256 _endTime) {
        timeLimitEnabled = _timeLimitEnabled;
        if (_timeLimitEnabled) {
            donationStartTime = _startTime == 0 ? block.timestamp : _startTime;
            donationEndTime = _endTime;
            require(_endTime == 0 || _endTime > donationStartTime, "End time must be after start time");
        }
        _owner = msg.sender;
    }
    

    receive() external payable {
        donate();
    }
    

    function donate() public payable {
        require(msg.value > 0, "Donation amount must be greater than 0");
        
        if (timeLimitEnabled) {
            require(block.timestamp >= donationStartTime, "Donation period has not started yet");
            require(donationEndTime == 0 || block.timestamp <= donationEndTime, "Donation period has ended");
        }
        
        // 如果是第一次捐赠，添加到捐赠者列表
        if (_donations[msg.sender] == 0) {
            _donors.push(msg.sender);
        }
        
        // 记录捐赠金额
        _donations[msg.sender] += msg.value;
        
        // 更新总捐赠金额
        totalDonations += msg.value;
        
        emit Donation(msg.sender, msg.value, block.timestamp);
    }
    

    function withdraw() public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");
        
        (bool success, ) = payable(_owner).call{value: balance}("");
        require(success, "Transfer failed");

    }
    
    function getDonation(address donor) public view returns (uint256) {
        return _donations[donor];
    }
    

    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }
    

    function getDonorCount() public view returns (uint256) {
        return _donors.length;
    }
    
    //获取捐赠金额最多的前 N 个地址
    function getTopDonors(uint256 topN) public view returns (address[] memory topDonors, uint256[] memory amounts) {
        require(topN > 0, "topN must be greater than 0");
        

        address[] memory donors = new address[](_donors.length);
        uint256[] memory donationAmounts = new uint256[](_donors.length);
        

        for (uint256 i = 0; i < _donors.length; i++) {
            donors[i] = _donors[i];
            donationAmounts[i] = _donations[_donors[i]];
        }
        
      
        for (uint256 i = 0; i < donors.length; i++) {
            for (uint256 j = 0; j < donors.length - i - 1; j++) {
                if (donationAmounts[j] < donationAmounts[j + 1]) {         
                    uint256 tempAmount = donationAmounts[j];
                    donationAmounts[j] = donationAmounts[j + 1];
                    donationAmounts[j + 1] = tempAmount;
                    address tempAddr = donors[j];
                    donors[j] = donors[j + 1];
                    donors[j + 1] = tempAddr;
                }
            }
        }
        
        uint256 returnCount = topN < donors.length ? topN : donors.length;
        
        topDonors = new address[](returnCount);
        amounts = new uint256[](returnCount);
        

        for (uint256 i = 0; i < returnCount; i++) {
            topDonors[i] = donors[i];
            amounts[i] = donationAmounts[i];
        }
    }
    
    function getTop3Donors() public view returns (address[] memory top3Donors, uint256[] memory top3Amounts) {
        return getTopDonors(3);
    }
    
    function setTimeLimit(bool _enabled, uint256 _startTime, uint256 _endTime) public onlyOwner {
        timeLimitEnabled = _enabled;
        if (_enabled) {
            donationStartTime = _startTime == 0 ? block.timestamp : _startTime;
            donationEndTime = _endTime;
            require(_endTime == 0 || _endTime > donationStartTime, "End time must be after start time");
        }
    }
    
    function isDonationPeriodActive() public view returns (bool) {
        if (!timeLimitEnabled) {
            return true;
        }
        if (block.timestamp < donationStartTime) {
            return false;
        }
        if (donationEndTime != 0 && block.timestamp > donationEndTime) {
            return false;
        }
        return true;
    }
    

    function getAllDonors() public view returns (address[] memory) {
        return _donors;
    }
}
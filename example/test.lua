Data = {
    count = {"",4},
    data = {"",1470}
}

function Decode(dataBytes)
    local nextBool = false
    if next(Data) then
        nextBool = true
    end
    local currentFieldKey = next(Data)
    local currentIndex = 1
    while nextBool do
        local fieldTable = Data[currentFieldKey]
        fieldTable[1] = string.sub(dataBytes,currentIndex,currentIndex+fieldTable[2]-1)
        print(currentIndex,currentIndex+fieldTable[2]-1,string.sub(fieldTable[1],1,10).."...")
        local fieldData = string.format("%X",fieldTable[1])
        print(currentFieldKey,fieldData)
        if next(Data,currentFieldKey) then
            nextBool = true
            currentFieldKey = next(Data,currentFieldKey)
            currentIndex = currentIndex + fieldTable[2]
        else
            nextBool = false
        end
    end
end

function Encode()
    local resultBytes = ""
    local nextBool = false
    if next(Data) then
        nextBool = true
    end
    local currentFieldKey = next(Data)
    local currentIndex = 1
    while nextBool do
        local fieldTable = Data[currentFieldKey]
        if fieldTable[2] == 1 then
            resultBytes = resultBytes..GenInt8()
        elseif fieldTable[2] == 2 then
            resultBytes = resultBytes..GenInt16()
        elseif fieldTable[2] == 4 then
            resultBytes = resultBytes..GenInt32()
        elseif fieldTable[2] > 4 then
            resultBytes = resultBytes..GenReadableBytes(fieldTable[2])
        end
        if next(Data,currentFieldKey) then
            nextBool = true
            currentFieldKey = next(Data,currentFieldKey)
            currentIndex = currentIndex + fieldTable[2]
        else
            nextBool = false
        end
    end
    return resultBytes
end


-- internal functions

-- 生成一个8位整数变量
function GenInt8()
    local num = math.random(0,255)
    local byte = string.char(num)
    return byte
end

-- 生成一个16位整数的字节变量
function GenInt16()
    local num = math.random(0, 65535) -- 生成0-65535之间的随机整数
    local byte1 = string.char(num % 256) -- 取低8位
    local byte2 = string.char(math.floor(num / 256)) -- 取高8位
    return byte2 .. byte1 -- 将两个字节拼接成一个字符串
end

-- 生成一个32位整数的字节变量
function GenInt32()
    local num = math.random(0, 4294967295) -- 生成0-4294967295之间的随机整数
    local byte1 = string.char(num % 256) -- 取低8位
    num = math.floor(num / 256)
    local byte2 = string.char(num % 256) -- 取次低8位
    num = math.floor(num / 256)
    local byte3 = string.char(num % 256) -- 取次高8位
    local byte4 = string.char(math.floor(num / 256)) -- 取高8位
    return byte4 .. byte3 .. byte2 .. byte1 -- 将四个字节拼接成一个字符串
end

-- 生成任意字节长度的可读字符字节变量
function GenReadableBytes(length)
    local chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    local bytes = ""
    for i = 1, length do
        local index = math.random(1, #chars)
        local char = string.sub(chars, index, index)
        bytes = bytes .. char
    end
    return bytes
end
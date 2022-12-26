Data = {
    count = {"",4},
    data = {"",1470}
}

function Parse(dataBytes)
    local nextBool = false
    if next(Data) then
        nextBool = true
    end
    local currentFieldKey = next(Data)
    local currentIndex = 1
    while nextBool do
        local fieldTable = Data[currentFieldKey]
        fieldTable[1] = string.sub(dataBytes,currentIndex,currentIndex+fieldTable[2]-1)
        print(currentIndex,currentIndex+fieldTable[2]-1,fieldTable[1])
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
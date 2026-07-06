function list()
    for i, v in ipairs(AllPlayers) do
        print(string.format("[%d] (%s) %s <%s>", i, v.userid, v.name, v.prefab))
    end
end

function GetPlayerData(kuid, uuid)
    local result = {
        uuid = uuid,
        kuid = kuid,
        success = false,
        error = nil,
        data = {}
    }
    
    -- 获取玩家
    local player = UserToPlayer(kuid)
    if not player then
        result.error = "Player not found with kuid: " .. tostring(kuid)
        return result
    end
    
    result.success = true
    
    -- 基础信息
    local data = {
        name = player.name or player.prefab or "Unknown",
        prefab = player.prefab,
        health = {
            current = player.components.health and player.components.health.currenthealth or 0,
            max = player.components.health and player.components.health.maxhealth or 0,
            percent = player.components.health and math.floor(player.components.health:GetPercent() * 100) or 0
        },
        hunger = {
            current = player.components.hunger and player.components.hunger.current or 0,
            max = player.components.hunger and player.components.hunger.max or 0,
            percent = player.components.hunger and math.floor(player.components.hunger:GetPercent() * 100) or 0
        },
        sanity = {
            current = player.components.sanity and player.components.sanity.current or 0,
            max = player.components.sanity and player.components.sanity.max or 0,
            percent = player.components.sanity and math.floor(player.components.sanity:GetPercent() * 100) or 0
        },
        temperature = player.components.temperature and math.floor(player.components.temperature.current) or 0,
        moisture = player.components.moisture and math.floor(player.components.moisture:GetMoisture()) or 0
    }
    
    local inv = player.components.inventory
    
    -- 手部装备
    local handItem = inv:GetEquippedItem(EQUIPSLOTS.HANDS)
    data.hand_equipment = handItem and {
        prefab = handItem.prefab,
        name = handItem.name or handItem.prefab,
        count = handItem.components.stackable and handItem.components.stackable.stacksize or 1,
        durability = handItem.components.finiteuses and math.floor(handItem.components.finiteuses:GetPercent() * 100) or nil,
        fuel = handItem.components.fueled and math.floor(handItem.components.fueled:GetPercent() * 100) or nil
    } or nil
    
    -- 头部装备
    local headItem = inv:GetEquippedItem(EQUIPSLOTS.HEAD)
    data.head_equipment = headItem and {
        prefab = headItem.prefab,
        name = headItem.name or headItem.prefab,
        count = headItem.components.stackable and headItem.components.stackable.stacksize or 1,
        durability = headItem.components.finiteuses and math.floor(headItem.components.finiteuses:GetPercent() * 100) or nil,
        armor = headItem.components.armor and math.floor(headItem.components.armor:GetPercent() * 100) or nil
    } or nil
    
    -- 身体装备（通常是背包或盔甲）
    local bodyItem = inv:GetEquippedItem(EQUIPSLOTS.BODY)
    local backpack = nil
    
    if bodyItem then
        -- 判断是否为背包（有容器组件）
        if bodyItem.components.container then
            backpack = bodyItem
            data.body_equipment = {
                type = "backpack",
                prefab = bodyItem.prefab,
                name = bodyItem.name or bodyItem.prefab,
                slots = bodyItem.components.container.numslots
            }
        else
            -- 普通盔甲/衣服
            data.body_equipment = {
                type = "armor",
                prefab = bodyItem.prefab,
                name = bodyItem.name or bodyItem.prefab,
                durability = bodyItem.components.armor and math.floor(bodyItem.components.armor:GetPercent() * 100) or nil
            }
        end
    else
        data.body_equipment = nil
    end
    
    -- 背包物品
    data.backpack_items = {}
    if backpack and backpack.components.container then
        local container = backpack.components.container
        for i = 1, container.numslots do
            local item = container.slots[i]
            if item then
                table.insert(data.backpack_items, {
                    slot = i,
                    prefab = item.prefab,
                    name = item.name or item.prefab,
                    count = item.components.stackable and item.components.stackable.stacksize or 1,
                    freshness = item.components.perishable and math.floor(item.components.perishable:GetPercent() * 100) or nil,
                    durability = item.components.finiteuses and math.floor(item.components.finiteuses:GetPercent() * 100) or nil
                })
            end
        end
    end
    
    -- 物品栏（15格）
    data.inventory_items = {}
    for i = 1, 15 do
        local item = inv.itemslots[i]
        if item then
            table.insert(data.inventory_items, {
                slot = i,
                prefab = item.prefab,
                name = item.name or item.prefab,
                count = item.components.stackable and item.components.stackable.stacksize or 1,
                freshness = item.components.perishable and math.floor(item.components.perishable:GetPercent() * 100) or nil,
                durability = item.components.finiteuses and math.floor(item.components.finiteuses:GetPercent() * 100) or nil,
                fuel = item.components.fueled and math.floor(item.components.fueled:GetPercent() * 100) or nil
            })
        end
    end
    
    -- 统计信息
    data.stats = {
        inventory_count = #data.inventory_items,
        backpack_count = #data.backpack_items,
        total_items = #data.inventory_items + #data.backpack_items
    }
    
    result.data = data
    return result
end

-- JSON 编码辅助函数（简单版）
function ToJSON(obj, indent)
    indent = indent or 0
    local spaces = string.rep("  ", indent)
    local result = {}
    
    if type(obj) == "table" then
        local isArray = #obj > 0
        table.insert(result, isArray and "[" or "{")
        
        local items = {}
        if isArray then
            for _, v in ipairs(obj) do
                table.insert(items, ToJSON(v, indent + 1))
            end
        else
            for k, v in pairs(obj) do
                local key = type(k) == "string" and ('"' .. k .. '":') or ('"' .. tostring(k) .. '":')
                table.insert(items, spaces .. "  " .. key .. ToJSON(v, indent + 1))
            end
        end
        
        table.insert(result, table.concat(items, ",\n"))
        table.insert(result, spaces .. (isArray and "]" or "}"))
    elseif type(obj) == "string" then
        return '"' .. obj:gsub('"', '\\"') .. '"'
    elseif type(obj) == "number" or type(obj) == "boolean" then
        return tostring(obj)
    elseif obj == nil then
        return "null"
    end
    
    return table.concat(result, "\n")
end


-- JSON 编码（压缩成一行）
function ToJSONStr(obj)
    local parts = {}
    
    local function serialize(o)
        if type(o) == "table" then
            local isArray = #o > 0
            table.insert(parts, isArray and "[" or "{")
            
            local first = true
            if isArray then
                for _, v in ipairs(o) do
                    if not first then table.insert(parts, ",") end
                    first = false
                    serialize(v)
                end
            else
                for k, v in pairs(o) do
                    if not first then table.insert(parts, ",") end
                    first = false
                    table.insert(parts, '"' .. tostring(k) .. '":')
                    serialize(v)
                end
            end
            
            table.insert(parts, isArray and "]" or "}")
        elseif type(o) == "string" then
            -- 转义特殊字符
            local s = o:gsub('\\', '\\\\'):gsub('"', '\\"'):gsub('\n', '\\n'):gsub('\r', '\\r'):gsub('\t', '\\t')
            table.insert(parts, '"' .. s .. '"')
        elseif type(o) == "number" then
            table.insert(parts, tostring(o))
        elseif type(o) == "boolean" then
            table.insert(parts, o and "true" or "false")
        elseif o == nil then
            table.insert(parts, "null")
        end
    end
    
    serialize(obj)
    return table.concat(parts)
end

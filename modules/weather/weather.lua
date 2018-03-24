local http = require("http")
local json = require("json")
local url = require("url")

local iconMap = {
	["clear-day"] = "wi-day-sunny",
	["clear-night"] = "wi-night-clear",
	["partly-cloudy-day"] = "wi-day-cloudy",
	["partly-cloudy-night"] = "wi-night-partly-cloudy"
}

function onRender(vars)
	if vars:get('latitude') == '' or vars:get('longitude') == '' then
		return 'error'
	end

	local response, error_message = http.get("https://api.darksky.net/forecast/" .. config.apiKey .. "/" .. vars:get('latitude') .. "," .. vars:get('longitude'))

	if error_message then
		error(error_message)
	end

	local res = json.decode(response.body)

	if iconMap[res.currently.icon] ~= nil then
		res.currently.icon = iconMap[res.currently.icon]
	else
		res.currently.icon = 'wi-' .. res.currently.icon
	end

	res.location = vars:get('location')

	return ctx:render('index', res)
end

function onGeocoding(vars)
	local query = url.build_query_string({
		["address"] = vars:get("location"),
		["key"] = config.geocodingKey
	})

	local response, error_message = http.get("https://maps.googleapis.com/maps/api/geocode/json?" .. query)

	if error_message then error(error_message) end

	local res = json.decode(response.body)

	if res.status ~= "OK" then
		error("Unable to get location")
	end

	return res
end
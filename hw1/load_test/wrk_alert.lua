request = function()
  local url = "/alert"
  return wrk.format("GET", url)
end
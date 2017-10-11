Feature: validate that the service is up and running

  Scenario: do a ping-pong request
    Given one environment variables "PAASTEK_CORE_ADDR" that matches "addr:port"
    And another variable PAASTEK_CORE_SCHEME=http
    And a service listening on PAASTEK_CORE_ADDR
    When I execute "curl ${PAASTEK_CORE_SCHEME}://${PAASTEK_CORE_ADDR}/ping"
    Then it returns "pong"


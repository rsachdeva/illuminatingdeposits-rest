# Ready to accept requests with db ready check
curl --location --request POST 'https://localhost:3000/v1/health' \
--header 'Content-Type: application/json'

###
#
#curl --location --request POST 'http://localhost:3000/v1/users' \
#--header 'Authorization: Basic cm9oaXRzM0Bkcmlubm92YXRpb25zLnVzOnBoYXNl' \
#POST http://localhost:3000/v1/users
#Authorization: Basic cm9oaXRzNDAwMDAwMDBAZHJpbm5vdmF0aW9ucy51czpwaGFzZQ==
#
####
#
#
####
#
# curl --location --request POST 'http://localhost:3000/v1/interests' \
#--header 'Content-Type: application/json' \
#--data-raw '{
#  "banks": [
#    {
#      "name": "HAPPIEST",
#      "deposits": [
#        {
#          "account": "1234",
#          "annualType": "Checking",
#          "annualRate%": 0,
#          "years": 1,
#          "amount": 1
#        },
#        {
#          "account": "1256",
#          "annualType": "CD",
#          "annualRate%": 24,
#          "years": 2,
#          "amount": 7700
#        },
#        {
#          "account": "1111",
#          "annualType": "CD",
#          "annualRate%": 1.01,
#          "years": 10,
#          "amount": 27000
#        }
#      ]
#    },
#    {
#      "name": "NICE",
#      "deposits": [
#        {
#          "account": "1234",
#          "annualType": "Brokered CD",
#          "annualRate%": 2.4,
#          "years": 7,
#          "amount": 10990
#        }
#      ]
#    },
#    {
#      "name": "ANGRY",
#      "deposits": [
#        {
#          "account": "1234",
#          "annualType": "Brokered CD",
#          "annualRate%": 5,
#          "years": 7,
#          "amount": 10990
#        },
#        {
#          "account": "9898",
#          "annualType": "CD",
#          "annualRate%": 2.22,
#          "years": 1,
#          "amount": 5500
#        }
#      ]
#    }
#  ]
##}'
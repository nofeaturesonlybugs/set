package set_test

import (
	"encoding/json"
	"testing"

	"github.com/nofeaturesonlybugs/set"
)

var benchmarkMapperJson string = `
[
	{
		"id": 1,
		"created_time": "2021-04-07",
		"modified_time": "2021-11-30",
		"price": 8,
		"quantity": 3,
		"total": 3,
		"customer_id": 1,
		"customer_first": "Burton",
		"customer_last": "Pratt",
		"vendor_id": 1,
		"vendor_name": "Dui Cum Sociis Company",
		"vendor_description": "egestas ligula. Nullam feugiat placerat velit. Quisque varius.",
		"vendor_contact_id": 1,
		"vendor_contact_first": "Mariko",
		"vendor_contact_last": "Buckner"
	},
	{
		"id": 2,
		"created_time": "2020-05-11",
		"modified_time": "2020-04-05",
		"price": 8,
		"quantity": 3,
		"total": 5,
		"customer_id": 2,
		"customer_first": "Dakota",
		"customer_last": "Wheeler",
		"vendor_id": 2,
		"vendor_name": "Mi Institute",
		"vendor_description": "id, ante. Nunc mauris sapien, cursus in,",
		"vendor_contact_id": 2,
		"vendor_contact_first": "Alisa",
		"vendor_contact_last": "Harmon"
	},
	{
		"id": 3,
		"created_time": "2020-09-13",
		"modified_time": "2020-12-13",
		"price": 7,
		"quantity": 5,
		"total": 8,
		"customer_id": 3,
		"customer_first": "Brock",
		"customer_last": "Thompson",
		"vendor_id": 3,
		"vendor_name": "Per Conubia Nostra Associates",
		"vendor_description": "Sed",
		"vendor_contact_id": 3,
		"vendor_contact_first": "Nathan",
		"vendor_contact_last": "Miranda"
	},
	{
		"id": 4,
		"created_time": "2020-04-25",
		"modified_time": "2020-03-23",
		"price": 8,
		"quantity": 5,
		"total": 8,
		"customer_id": 4,
		"customer_first": "Gregory",
		"customer_last": "Gilmore",
		"vendor_id": 4,
		"vendor_name": "Aliquam LLC",
		"vendor_description": "Nunc ut erat. Sed",
		"vendor_contact_id": 4,
		"vendor_contact_first": "Serena",
		"vendor_contact_last": "Chambers"
	},
	{
		"id": 5,
		"created_time": "2020-04-27",
		"modified_time": "2020-08-21",
		"price": 1,
		"quantity": 7,
		"total": 8,
		"customer_id": 5,
		"customer_first": "Cooper",
		"customer_last": "Larson",
		"vendor_id": 5,
		"vendor_name": "Primis In Faucibus PC",
		"vendor_description": "Aenean sed pede nec ante blandit viverra. Donec",
		"vendor_contact_id": 5,
		"vendor_contact_first": "Colorado",
		"vendor_contact_last": "Rivera"
	},
	{
		"id": 6,
		"created_time": "2021-09-27",
		"modified_time": "2021-05-15",
		"price": 1,
		"quantity": 2,
		"total": 10,
		"customer_id": 6,
		"customer_first": "Alma",
		"customer_last": "Vincent",
		"vendor_id": 6,
		"vendor_name": "Ac Corp.",
		"vendor_description": "per",
		"vendor_contact_id": 6,
		"vendor_contact_first": "Hayden",
		"vendor_contact_last": "Santana"
	},
	{
		"id": 7,
		"created_time": "2020-10-12",
		"modified_time": "2020-07-30",
		"price": 7,
		"quantity": 9,
		"total": 6,
		"customer_id": 7,
		"customer_first": "Rebecca",
		"customer_last": "Cross",
		"vendor_id": 7,
		"vendor_name": "Ut Tincidunt Orci Consulting",
		"vendor_description": "at, velit. Pellentesque ultricies dignissim lacus. Aliquam rutrum",
		"vendor_contact_id": 7,
		"vendor_contact_first": "Lysandra",
		"vendor_contact_last": "Morton"
	},
	{
		"id": 8,
		"created_time": "2021-11-15",
		"modified_time": "2021-01-22",
		"price": 9,
		"quantity": 7,
		"total": 3,
		"customer_id": 8,
		"customer_first": "Troy",
		"customer_last": "Horton",
		"vendor_id": 8,
		"vendor_name": "Nullam Scelerisque Associates",
		"vendor_description": "nibh enim, gravida sit amet, dapibus id, blandit",
		"vendor_contact_id": 8,
		"vendor_contact_first": "Keefe",
		"vendor_contact_last": "Thomas"
	},
	{
		"id": 9,
		"created_time": "2020-05-14",
		"modified_time": "2020-06-21",
		"price": 10,
		"quantity": 1,
		"total": 1,
		"customer_id": 9,
		"customer_first": "Xantha",
		"customer_last": "Dawson",
		"vendor_id": 9,
		"vendor_name": "Nec Luctus Inc.",
		"vendor_description": "facilisis eget, ipsum. Donec sollicitudin adipiscing ligula. Aenean gravida",
		"vendor_contact_id": 9,
		"vendor_contact_first": "Bruce",
		"vendor_contact_last": "Roman"
	},
	{
		"id": 10,
		"created_time": "2021-09-20",
		"modified_time": "2020-05-11",
		"price": 8,
		"quantity": 5,
		"total": 9,
		"customer_id": 10,
		"customer_first": "Kendall",
		"customer_last": "Odonnell",
		"vendor_id": 10,
		"vendor_name": "Dictum Sapien PC",
		"vendor_description": "odio vel est tempor",
		"vendor_contact_id": 10,
		"vendor_contact_first": "Winifred",
		"vendor_contact_last": "Conrad"
	},
	{
		"id": 11,
		"created_time": "2020-08-04",
		"modified_time": "2022-01-24",
		"price": 8,
		"quantity": 2,
		"total": 5,
		"customer_id": 11,
		"customer_first": "Erin",
		"customer_last": "Mccormick",
		"vendor_id": 11,
		"vendor_name": "Maecenas Limited",
		"vendor_description": "congue a, aliquet vel, vulputate eu, odio. Phasellus",
		"vendor_contact_id": 11,
		"vendor_contact_first": "Sawyer",
		"vendor_contact_last": "Williamson"
	},
	{
		"id": 12,
		"created_time": "2021-12-01",
		"modified_time": "2020-06-10",
		"price": 6,
		"quantity": 10,
		"total": 7,
		"customer_id": 12,
		"customer_first": "Dahlia",
		"customer_last": "Whitley",
		"vendor_id": 12,
		"vendor_name": "Auctor Velit Aliquam Inc.",
		"vendor_description": "gravida mauris ut mi. Duis risus odio, auctor",
		"vendor_contact_id": 12,
		"vendor_contact_first": "Jorden",
		"vendor_contact_last": "Douglas"
	},
	{
		"id": 13,
		"created_time": "2021-09-12",
		"modified_time": "2021-05-18",
		"price": 9,
		"quantity": 6,
		"total": 3,
		"customer_id": 13,
		"customer_first": "Marshall",
		"customer_last": "Benton",
		"vendor_id": 13,
		"vendor_name": "Ante Bibendum Corp.",
		"vendor_description": "ante. Nunc mauris sapien, cursus in, hendrerit",
		"vendor_contact_id": 13,
		"vendor_contact_first": "Jocelyn",
		"vendor_contact_last": "Moreno"
	},
	{
		"id": 14,
		"created_time": "2021-05-06",
		"modified_time": "2021-12-30",
		"price": 6,
		"quantity": 10,
		"total": 4,
		"customer_id": 14,
		"customer_first": "Derek",
		"customer_last": "Ward",
		"vendor_id": 14,
		"vendor_name": "Tristique Pellentesque Company",
		"vendor_description": "orci, in consequat enim diam vel arcu. Curabitur ut",
		"vendor_contact_id": 14,
		"vendor_contact_first": "India",
		"vendor_contact_last": "Vaughn"
	},
	{
		"id": 15,
		"created_time": "2021-11-29",
		"modified_time": "2021-01-21",
		"price": 9,
		"quantity": 9,
		"total": 9,
		"customer_id": 15,
		"customer_first": "Tamekah",
		"customer_last": "Mcneil",
		"vendor_id": 15,
		"vendor_name": "Phasellus Dapibus Quam Associates",
		"vendor_description": "dis parturient montes, nascetur",
		"vendor_contact_id": 15,
		"vendor_contact_first": "Jordan",
		"vendor_contact_last": "Levy"
	},
	{
		"id": 16,
		"created_time": "2020-07-05",
		"modified_time": "2021-12-05",
		"price": 9,
		"quantity": 7,
		"total": 6,
		"customer_id": 16,
		"customer_first": "Uma",
		"customer_last": "Olson",
		"vendor_id": 16,
		"vendor_name": "Aliquam PC",
		"vendor_description": "mollis. Phasellus libero mauris, aliquam eu, accumsan sed, facilisis",
		"vendor_contact_id": 16,
		"vendor_contact_first": "Fiona",
		"vendor_contact_last": "Mclean"
	},
	{
		"id": 17,
		"created_time": "2021-03-04",
		"modified_time": "2021-12-23",
		"price": 3,
		"quantity": 10,
		"total": 3,
		"customer_id": 17,
		"customer_first": "Barbara",
		"customer_last": "Roberson",
		"vendor_id": 17,
		"vendor_name": "Vivamus Sit LLP",
		"vendor_description": "Etiam imperdiet dictum magna. Ut",
		"vendor_contact_id": 17,
		"vendor_contact_first": "Hop",
		"vendor_contact_last": "Vance"
	},
	{
		"id": 18,
		"created_time": "2020-09-22",
		"modified_time": "2021-05-05",
		"price": 5,
		"quantity": 2,
		"total": 7,
		"customer_id": 18,
		"customer_first": "Lenore",
		"customer_last": "Fitzgerald",
		"vendor_id": 18,
		"vendor_name": "Lacinia Orci Consectetuer Industries",
		"vendor_description": "aliquet molestie tellus. Aenean egestas hendrerit neque. In ornare sagittis",
		"vendor_contact_id": 18,
		"vendor_contact_first": "Tobias",
		"vendor_contact_last": "Dunlap"
	},
	{
		"id": 19,
		"created_time": "2020-03-05",
		"modified_time": "2020-05-26",
		"price": 4,
		"quantity": 6,
		"total": 5,
		"customer_id": 19,
		"customer_first": "Wylie",
		"customer_last": "Cook",
		"vendor_id": 19,
		"vendor_name": "Donec Consulting",
		"vendor_description": "adipiscing.",
		"vendor_contact_id": 19,
		"vendor_contact_first": "Abdul",
		"vendor_contact_last": "Anderson"
	},
	{
		"id": 20,
		"created_time": "2022-01-11",
		"modified_time": "2020-12-29",
		"price": 2,
		"quantity": 5,
		"total": 6,
		"customer_id": 20,
		"customer_first": "Ross",
		"customer_last": "Wilson",
		"vendor_id": 20,
		"vendor_name": "Sociis Natoque Penatibus LLP",
		"vendor_description": "Suspendisse commodo tincidunt nibh. Phasellus",
		"vendor_contact_id": 20,
		"vendor_contact_first": "Libby",
		"vendor_contact_last": "Gould"
	},
	{
		"id": 21,
		"created_time": "2022-01-21",
		"modified_time": "2021-10-29",
		"price": 7,
		"quantity": 9,
		"total": 8,
		"customer_id": 21,
		"customer_first": "Ray",
		"customer_last": "Vazquez",
		"vendor_id": 21,
		"vendor_name": "Imperdiet Non Vestibulum Inc.",
		"vendor_description": "In condimentum. Donec at",
		"vendor_contact_id": 21,
		"vendor_contact_first": "Gail",
		"vendor_contact_last": "Bean"
	},
	{
		"id": 22,
		"created_time": "2020-05-07",
		"modified_time": "2021-02-11",
		"price": 5,
		"quantity": 8,
		"total": 1,
		"customer_id": 22,
		"customer_first": "Xyla",
		"customer_last": "Townsend",
		"vendor_id": 22,
		"vendor_name": "Et Industries",
		"vendor_description": "elit, pretium et, rutrum non, hendrerit id,",
		"vendor_contact_id": 22,
		"vendor_contact_first": "Kamal",
		"vendor_contact_last": "Gonzalez"
	},
	{
		"id": 23,
		"created_time": "2020-08-25",
		"modified_time": "2021-01-25",
		"price": 6,
		"quantity": 4,
		"total": 1,
		"customer_id": 23,
		"customer_first": "Venus",
		"customer_last": "Atkins",
		"vendor_id": 23,
		"vendor_name": "Sem Limited",
		"vendor_description": "Sed pharetra, felis eget varius ultrices, mauris ipsum",
		"vendor_contact_id": 23,
		"vendor_contact_first": "Carly",
		"vendor_contact_last": "Travis"
	},
	{
		"id": 24,
		"created_time": "2020-06-05",
		"modified_time": "2021-05-20",
		"price": 5,
		"quantity": 10,
		"total": 1,
		"customer_id": 24,
		"customer_first": "Russell",
		"customer_last": "Cunningham",
		"vendor_id": 24,
		"vendor_name": "Bibendum Donec Corporation",
		"vendor_description": "odio. Etiam ligula tortor, dictum eu, placerat eget, venenatis a,",
		"vendor_contact_id": 24,
		"vendor_contact_first": "Lacey",
		"vendor_contact_last": "Leonard"
	},
	{
		"id": 25,
		"created_time": "2020-10-16",
		"modified_time": "2021-01-31",
		"price": 10,
		"quantity": 2,
		"total": 7,
		"customer_id": 25,
		"customer_first": "Jennifer",
		"customer_last": "Douglas",
		"vendor_id": 25,
		"vendor_name": "Non Sapien Inc.",
		"vendor_description": "rutrum non, hendrerit id, ante. Nunc mauris sapien, cursus",
		"vendor_contact_id": 25,
		"vendor_contact_first": "Hashim",
		"vendor_contact_last": "Bernard"
	},
	{
		"id": 26,
		"created_time": "2020-07-31",
		"modified_time": "2022-01-30",
		"price": 6,
		"quantity": 6,
		"total": 8,
		"customer_id": 26,
		"customer_first": "Catherine",
		"customer_last": "Cooke",
		"vendor_id": 26,
		"vendor_name": "Nonummy Ipsum Non PC",
		"vendor_description": "magnis dis parturient montes, nascetur ridiculus",
		"vendor_contact_id": 26,
		"vendor_contact_first": "Warren",
		"vendor_contact_last": "Stark"
	},
	{
		"id": 27,
		"created_time": "2021-10-01",
		"modified_time": "2021-09-17",
		"price": 1,
		"quantity": 9,
		"total": 10,
		"customer_id": 27,
		"customer_first": "Shafira",
		"customer_last": "Dale",
		"vendor_id": 27,
		"vendor_name": "Ornare Elit Elit Associates",
		"vendor_description": "nonummy ipsum non arcu. Vivamus sit amet",
		"vendor_contact_id": 27,
		"vendor_contact_first": "Victoria",
		"vendor_contact_last": "Joyner"
	},
	{
		"id": 28,
		"created_time": "2021-10-03",
		"modified_time": "2020-07-03",
		"price": 2,
		"quantity": 8,
		"total": 2,
		"customer_id": 28,
		"customer_first": "Violet",
		"customer_last": "Bradford",
		"vendor_id": 28,
		"vendor_name": "Posuere At LLC",
		"vendor_description": "venenatis a, magna. Lorem",
		"vendor_contact_id": 28,
		"vendor_contact_first": "Drake",
		"vendor_contact_last": "Mendez"
	},
	{
		"id": 29,
		"created_time": "2021-05-14",
		"modified_time": "2020-12-07",
		"price": 5,
		"quantity": 10,
		"total": 7,
		"customer_id": 29,
		"customer_first": "Daniel",
		"customer_last": "Parker",
		"vendor_id": 29,
		"vendor_name": "Nullam Ut Nisi Corp.",
		"vendor_description": "mus. Aenean eget magna.",
		"vendor_contact_id": 29,
		"vendor_contact_first": "Kuame",
		"vendor_contact_last": "Gilliam"
	},
	{
		"id": 30,
		"created_time": "2020-07-19",
		"modified_time": "2020-11-09",
		"price": 8,
		"quantity": 4,
		"total": 6,
		"customer_id": 30,
		"customer_first": "Mark",
		"customer_last": "Berg",
		"vendor_id": 30,
		"vendor_name": "Lacus Pede PC",
		"vendor_description": "ac urna. Ut tincidunt vehicula",
		"vendor_contact_id": 30,
		"vendor_contact_first": "Halla",
		"vendor_contact_last": "Yates"
	},
	{
		"id": 31,
		"created_time": "2020-03-18",
		"modified_time": "2020-03-15",
		"price": 5,
		"quantity": 6,
		"total": 8,
		"customer_id": 31,
		"customer_first": "Candace",
		"customer_last": "Patel",
		"vendor_id": 31,
		"vendor_name": "Cras Eu Tellus LLP",
		"vendor_description": "id, ante. Nunc mauris sapien, cursus in, hendrerit consectetuer, cursus",
		"vendor_contact_id": 31,
		"vendor_contact_first": "Ivor",
		"vendor_contact_last": "Maddox"
	},
	{
		"id": 32,
		"created_time": "2021-03-25",
		"modified_time": "2020-06-04",
		"price": 1,
		"quantity": 2,
		"total": 10,
		"customer_id": 32,
		"customer_first": "Iris",
		"customer_last": "Larson",
		"vendor_id": 32,
		"vendor_name": "Euismod Enim PC",
		"vendor_description": "massa non",
		"vendor_contact_id": 32,
		"vendor_contact_first": "Victoria",
		"vendor_contact_last": "Hood"
	},
	{
		"id": 33,
		"created_time": "2021-08-06",
		"modified_time": "2021-11-10",
		"price": 7,
		"quantity": 8,
		"total": 9,
		"customer_id": 33,
		"customer_first": "Caesar",
		"customer_last": "Roy",
		"vendor_id": 33,
		"vendor_name": "Ut Tincidunt Foundation",
		"vendor_description": "magna a neque. Nullam ut nisi a odio semper",
		"vendor_contact_id": 33,
		"vendor_contact_first": "Brendan",
		"vendor_contact_last": "Spence"
	},
	{
		"id": 34,
		"created_time": "2021-08-14",
		"modified_time": "2020-07-16",
		"price": 5,
		"quantity": 4,
		"total": 9,
		"customer_id": 34,
		"customer_first": "Gemma",
		"customer_last": "Mueller",
		"vendor_id": 34,
		"vendor_name": "Ut Semper LLC",
		"vendor_description": "lectus pede et",
		"vendor_contact_id": 34,
		"vendor_contact_first": "Dustin",
		"vendor_contact_last": "Higgins"
	},
	{
		"id": 35,
		"created_time": "2020-02-28",
		"modified_time": "2021-06-05",
		"price": 6,
		"quantity": 6,
		"total": 6,
		"customer_id": 35,
		"customer_first": "Bell",
		"customer_last": "Burgess",
		"vendor_id": 35,
		"vendor_name": "Proin Vel Nisl Foundation",
		"vendor_description": "Quisque fringilla euismod enim. Etiam",
		"vendor_contact_id": 35,
		"vendor_contact_first": "Nayda",
		"vendor_contact_last": "Jacobs"
	},
	{
		"id": 36,
		"created_time": "2021-08-13",
		"modified_time": "2021-07-20",
		"price": 2,
		"quantity": 7,
		"total": 6,
		"customer_id": 36,
		"customer_first": "Martin",
		"customer_last": "Wilkerson",
		"vendor_id": 36,
		"vendor_name": "Rhoncus Limited",
		"vendor_description": "at pretium aliquet, metus urna convallis erat, eget tincidunt",
		"vendor_contact_id": 36,
		"vendor_contact_first": "Amber",
		"vendor_contact_last": "York"
	},
	{
		"id": 37,
		"created_time": "2021-04-03",
		"modified_time": "2021-05-05",
		"price": 6,
		"quantity": 8,
		"total": 2,
		"customer_id": 37,
		"customer_first": "Hasad",
		"customer_last": "Riddle",
		"vendor_id": 37,
		"vendor_name": "A Feugiat Tellus LLC",
		"vendor_description": "at risus. Nunc",
		"vendor_contact_id": 37,
		"vendor_contact_first": "Karina",
		"vendor_contact_last": "Brennan"
	},
	{
		"id": 38,
		"created_time": "2021-05-11",
		"modified_time": "2020-03-11",
		"price": 7,
		"quantity": 10,
		"total": 9,
		"customer_id": 38,
		"customer_first": "Sydnee",
		"customer_last": "Gregory",
		"vendor_id": 38,
		"vendor_name": "Consectetuer Rhoncus Nullam Ltd",
		"vendor_description": "semper",
		"vendor_contact_id": 38,
		"vendor_contact_first": "Hedwig",
		"vendor_contact_last": "Watson"
	},
	{
		"id": 39,
		"created_time": "2021-01-31",
		"modified_time": "2020-12-26",
		"price": 5,
		"quantity": 3,
		"total": 10,
		"customer_id": 39,
		"customer_first": "Colin",
		"customer_last": "Lott",
		"vendor_id": 39,
		"vendor_name": "Elit Sed PC",
		"vendor_description": "nec, leo. Morbi neque tellus,",
		"vendor_contact_id": 39,
		"vendor_contact_first": "Burke",
		"vendor_contact_last": "David"
	},
	{
		"id": 40,
		"created_time": "2021-11-05",
		"modified_time": "2021-08-05",
		"price": 2,
		"quantity": 3,
		"total": 7,
		"customer_id": 40,
		"customer_first": "Cyrus",
		"customer_last": "Parker",
		"vendor_id": 40,
		"vendor_name": "Mauris Ut Foundation",
		"vendor_description": "eu, odio. Phasellus at augue id ante dictum",
		"vendor_contact_id": 40,
		"vendor_contact_first": "Macon",
		"vendor_contact_last": "Hewitt"
	},
	{
		"id": 41,
		"created_time": "2021-03-08",
		"modified_time": "2021-05-07",
		"price": 2,
		"quantity": 10,
		"total": 8,
		"customer_id": 41,
		"customer_first": "Raja",
		"customer_last": "Merritt",
		"vendor_id": 41,
		"vendor_name": "Nunc Est Mollis Company",
		"vendor_description": "vitae, sodales at, velit. Pellentesque ultricies dignissim",
		"vendor_contact_id": 41,
		"vendor_contact_first": "Palmer",
		"vendor_contact_last": "Mckay"
	},
	{
		"id": 42,
		"created_time": "2021-06-06",
		"modified_time": "2021-09-25",
		"price": 9,
		"quantity": 1,
		"total": 8,
		"customer_id": 42,
		"customer_first": "Melanie",
		"customer_last": "Rowe",
		"vendor_id": 42,
		"vendor_name": "Lacus Varius Et Corporation",
		"vendor_description": "tellus, imperdiet",
		"vendor_contact_id": 42,
		"vendor_contact_first": "Fletcher",
		"vendor_contact_last": "Irwin"
	},
	{
		"id": 43,
		"created_time": "2021-02-27",
		"modified_time": "2021-08-26",
		"price": 7,
		"quantity": 3,
		"total": 9,
		"customer_id": 43,
		"customer_first": "Remedios",
		"customer_last": "Wilder",
		"vendor_id": 43,
		"vendor_name": "Non PC",
		"vendor_description": "mus. Aenean eget",
		"vendor_contact_id": 43,
		"vendor_contact_first": "Haviva",
		"vendor_contact_last": "Contreras"
	},
	{
		"id": 44,
		"created_time": "2021-07-11",
		"modified_time": "2021-11-29",
		"price": 8,
		"quantity": 4,
		"total": 10,
		"customer_id": 44,
		"customer_first": "Carissa",
		"customer_last": "Nieves",
		"vendor_id": 44,
		"vendor_name": "In Corp.",
		"vendor_description": "vel, convallis in, cursus et,",
		"vendor_contact_id": 44,
		"vendor_contact_first": "Abraham",
		"vendor_contact_last": "Sanford"
	},
	{
		"id": 45,
		"created_time": "2021-12-19",
		"modified_time": "2021-06-30",
		"price": 1,
		"quantity": 4,
		"total": 1,
		"customer_id": 45,
		"customer_first": "Inez",
		"customer_last": "Romero",
		"vendor_id": 45,
		"vendor_name": "Lorem Ipsum Industries",
		"vendor_description": "fringilla euismod enim. Etiam",
		"vendor_contact_id": 45,
		"vendor_contact_first": "Price",
		"vendor_contact_last": "Alston"
	},
	{
		"id": 46,
		"created_time": "2021-07-08",
		"modified_time": "2021-11-05",
		"price": 3,
		"quantity": 7,
		"total": 3,
		"customer_id": 46,
		"customer_first": "Jack",
		"customer_last": "Rhodes",
		"vendor_id": 46,
		"vendor_name": "Odio LLC",
		"vendor_description": "vestibulum massa rutrum magna. Cras convallis convallis dolor.",
		"vendor_contact_id": 46,
		"vendor_contact_first": "Emerson",
		"vendor_contact_last": "Acevedo"
	},
	{
		"id": 47,
		"created_time": "2020-10-02",
		"modified_time": "2021-01-10",
		"price": 4,
		"quantity": 10,
		"total": 9,
		"customer_id": 47,
		"customer_first": "Avram",
		"customer_last": "Pearson",
		"vendor_id": 47,
		"vendor_name": "Ornare Tortor PC",
		"vendor_description": "sagittis. Duis gravida. Praesent eu nulla at",
		"vendor_contact_id": 47,
		"vendor_contact_first": "Stuart",
		"vendor_contact_last": "Hickman"
	},
	{
		"id": 48,
		"created_time": "2022-02-03",
		"modified_time": "2022-01-06",
		"price": 5,
		"quantity": 8,
		"total": 8,
		"customer_id": 48,
		"customer_first": "Nolan",
		"customer_last": "Dawson",
		"vendor_id": 48,
		"vendor_name": "A Company",
		"vendor_description": "Cras",
		"vendor_contact_id": 48,
		"vendor_contact_first": "Wilma",
		"vendor_contact_last": "Reese"
	},
	{
		"id": 49,
		"created_time": "2020-05-07",
		"modified_time": "2020-04-02",
		"price": 4,
		"quantity": 8,
		"total": 2,
		"customer_id": 49,
		"customer_first": "Erica",
		"customer_last": "Buckley",
		"vendor_id": 49,
		"vendor_name": "Sagittis Semper Institute",
		"vendor_description": "sodales at, velit. Pellentesque ultricies",
		"vendor_contact_id": 49,
		"vendor_contact_first": "Colby",
		"vendor_contact_last": "Frederick"
	},
	{
		"id": 50,
		"created_time": "2022-02-08",
		"modified_time": "2021-10-27",
		"price": 6,
		"quantity": 10,
		"total": 7,
		"customer_id": 50,
		"customer_first": "Mason",
		"customer_last": "Conway",
		"vendor_id": 50,
		"vendor_name": "Posuere Cubilia Curae; Foundation",
		"vendor_description": "lorem ut aliquam iaculis, lacus pede sagittis augue,",
		"vendor_contact_id": 50,
		"vendor_contact_first": "Emerald",
		"vendor_contact_last": "Berry"
	},
	{
		"id": 51,
		"created_time": "2021-09-24",
		"modified_time": "2021-04-16",
		"price": 10,
		"quantity": 3,
		"total": 4,
		"customer_id": 51,
		"customer_first": "Macey",
		"customer_last": "Larsen",
		"vendor_id": 51,
		"vendor_name": "Cum Corporation",
		"vendor_description": "et pede. Nunc sed orci lobortis augue",
		"vendor_contact_id": 51,
		"vendor_contact_first": "Delilah",
		"vendor_contact_last": "Norton"
	},
	{
		"id": 52,
		"created_time": "2020-03-17",
		"modified_time": "2021-02-11",
		"price": 9,
		"quantity": 3,
		"total": 9,
		"customer_id": 52,
		"customer_first": "Meghan",
		"customer_last": "Jones",
		"vendor_id": 52,
		"vendor_name": "Dictum Mi PC",
		"vendor_description": "tristique",
		"vendor_contact_id": 52,
		"vendor_contact_first": "Dieter",
		"vendor_contact_last": "Williams"
	},
	{
		"id": 53,
		"created_time": "2021-05-06",
		"modified_time": "2022-02-08",
		"price": 4,
		"quantity": 1,
		"total": 2,
		"customer_id": 53,
		"customer_first": "Orlando",
		"customer_last": "Huber",
		"vendor_id": 53,
		"vendor_name": "Neque LLP",
		"vendor_description": "pharetra ut, pharetra sed, hendrerit a, arcu. Sed",
		"vendor_contact_id": 53,
		"vendor_contact_first": "Ivor",
		"vendor_contact_last": "Foley"
	},
	{
		"id": 54,
		"created_time": "2021-08-08",
		"modified_time": "2021-09-02",
		"price": 1,
		"quantity": 6,
		"total": 3,
		"customer_id": 54,
		"customer_first": "Macy",
		"customer_last": "Calderon",
		"vendor_id": 54,
		"vendor_name": "Egestas Fusce Aliquet Associates",
		"vendor_description": "at, velit. Pellentesque ultricies",
		"vendor_contact_id": 54,
		"vendor_contact_first": "Elvis",
		"vendor_contact_last": "Lucas"
	},
	{
		"id": 55,
		"created_time": "2021-08-16",
		"modified_time": "2020-02-16",
		"price": 8,
		"quantity": 6,
		"total": 4,
		"customer_id": 55,
		"customer_first": "Phyllis",
		"customer_last": "Snider",
		"vendor_id": 55,
		"vendor_name": "Lorem Donec Elementum Foundation",
		"vendor_description": "Mauris eu turpis. Nulla aliquet. Proin velit. Sed malesuada",
		"vendor_contact_id": 55,
		"vendor_contact_first": "Kaden",
		"vendor_contact_last": "Snyder"
	},
	{
		"id": 56,
		"created_time": "2020-08-18",
		"modified_time": "2021-03-07",
		"price": 3,
		"quantity": 8,
		"total": 5,
		"customer_id": 56,
		"customer_first": "Kiona",
		"customer_last": "Weber",
		"vendor_id": 56,
		"vendor_name": "Diam At Pretium LLC",
		"vendor_description": "ipsum. Phasellus vitae",
		"vendor_contact_id": 56,
		"vendor_contact_first": "Philip",
		"vendor_contact_last": "Brennan"
	},
	{
		"id": 57,
		"created_time": "2021-04-16",
		"modified_time": "2020-06-27",
		"price": 2,
		"quantity": 6,
		"total": 8,
		"customer_id": 57,
		"customer_first": "Maia",
		"customer_last": "Morin",
		"vendor_id": 57,
		"vendor_name": "Diam At Company",
		"vendor_description": "cursus. Nunc mauris elit, dictum",
		"vendor_contact_id": 57,
		"vendor_contact_first": "Jessica",
		"vendor_contact_last": "Henry"
	},
	{
		"id": 58,
		"created_time": "2020-10-18",
		"modified_time": "2022-01-30",
		"price": 10,
		"quantity": 4,
		"total": 3,
		"customer_id": 58,
		"customer_first": "Mira",
		"customer_last": "Wong",
		"vendor_id": 58,
		"vendor_name": "Est Industries",
		"vendor_description": "tincidunt vehicula risus. Nulla eget metus eu",
		"vendor_contact_id": 58,
		"vendor_contact_first": "Iris",
		"vendor_contact_last": "Cobb"
	},
	{
		"id": 59,
		"created_time": "2020-09-19",
		"modified_time": "2021-04-04",
		"price": 4,
		"quantity": 7,
		"total": 1,
		"customer_id": 59,
		"customer_first": "Shana",
		"customer_last": "Stein",
		"vendor_id": 59,
		"vendor_name": "Mauris Rhoncus Limited",
		"vendor_description": "id magna et ipsum cursus vestibulum. Mauris magna. Duis",
		"vendor_contact_id": 59,
		"vendor_contact_first": "Kaye",
		"vendor_contact_last": "Macias"
	},
	{
		"id": 60,
		"created_time": "2021-10-28",
		"modified_time": "2020-06-05",
		"price": 2,
		"quantity": 3,
		"total": 4,
		"customer_id": 60,
		"customer_first": "Nina",
		"customer_last": "Wyatt",
		"vendor_id": 60,
		"vendor_name": "Quis Turpis Limited",
		"vendor_description": "venenatis a,",
		"vendor_contact_id": 60,
		"vendor_contact_first": "Shaine",
		"vendor_contact_last": "Conner"
	},
	{
		"id": 61,
		"created_time": "2022-01-06",
		"modified_time": "2021-07-14",
		"price": 6,
		"quantity": 6,
		"total": 4,
		"customer_id": 61,
		"customer_first": "Cynthia",
		"customer_last": "Wong",
		"vendor_id": 61,
		"vendor_name": "Magnis Incorporated",
		"vendor_description": "nunc risus varius orci, in",
		"vendor_contact_id": 61,
		"vendor_contact_first": "Iliana",
		"vendor_contact_last": "Norton"
	},
	{
		"id": 62,
		"created_time": "2021-09-18",
		"modified_time": "2020-08-29",
		"price": 6,
		"quantity": 1,
		"total": 1,
		"customer_id": 62,
		"customer_first": "Xandra",
		"customer_last": "Waller",
		"vendor_id": 62,
		"vendor_name": "Suscipit Inc.",
		"vendor_description": "Sed diam lorem, auctor quis, tristique ac, eleifend vitae, erat.",
		"vendor_contact_id": 62,
		"vendor_contact_first": "Akeem",
		"vendor_contact_last": "Velazquez"
	},
	{
		"id": 63,
		"created_time": "2021-04-01",
		"modified_time": "2020-03-10",
		"price": 5,
		"quantity": 3,
		"total": 8,
		"customer_id": 63,
		"customer_first": "Erich",
		"customer_last": "Cantrell",
		"vendor_id": 63,
		"vendor_name": "Mauris Corp.",
		"vendor_description": "sodales purus, in molestie tortor",
		"vendor_contact_id": 63,
		"vendor_contact_first": "Emerald",
		"vendor_contact_last": "Carson"
	},
	{
		"id": 64,
		"created_time": "2021-03-28",
		"modified_time": "2021-10-08",
		"price": 3,
		"quantity": 1,
		"total": 6,
		"customer_id": 64,
		"customer_first": "Oprah",
		"customer_last": "Frederick",
		"vendor_id": 64,
		"vendor_name": "Velit Cras Inc.",
		"vendor_description": "per conubia nostra, per inceptos hymenaeos. Mauris ut quam",
		"vendor_contact_id": 64,
		"vendor_contact_first": "Rebecca",
		"vendor_contact_last": "Francis"
	},
	{
		"id": 65,
		"created_time": "2022-01-24",
		"modified_time": "2021-06-17",
		"price": 9,
		"quantity": 5,
		"total": 10,
		"customer_id": 65,
		"customer_first": "Blossom",
		"customer_last": "Terrell",
		"vendor_id": 65,
		"vendor_name": "Faucibus Id Industries",
		"vendor_description": "et nunc. Quisque ornare tortor",
		"vendor_contact_id": 65,
		"vendor_contact_first": "Jack",
		"vendor_contact_last": "Weiss"
	},
	{
		"id": 66,
		"created_time": "2020-12-20",
		"modified_time": "2021-04-13",
		"price": 10,
		"quantity": 3,
		"total": 6,
		"customer_id": 66,
		"customer_first": "Lani",
		"customer_last": "Goff",
		"vendor_id": 66,
		"vendor_name": "Sed Institute",
		"vendor_description": "commodo hendrerit. Donec porttitor tellus",
		"vendor_contact_id": 66,
		"vendor_contact_first": "Meredith",
		"vendor_contact_last": "Merritt"
	},
	{
		"id": 67,
		"created_time": "2021-10-16",
		"modified_time": "2020-09-02",
		"price": 9,
		"quantity": 5,
		"total": 6,
		"customer_id": 67,
		"customer_first": "Charles",
		"customer_last": "Guzman",
		"vendor_id": 67,
		"vendor_name": "Donec Tempor Est Industries",
		"vendor_description": "Duis sit amet diam eu dolor",
		"vendor_contact_id": 67,
		"vendor_contact_first": "Kyle",
		"vendor_contact_last": "Saunders"
	},
	{
		"id": 68,
		"created_time": "2021-07-26",
		"modified_time": "2021-01-07",
		"price": 3,
		"quantity": 8,
		"total": 5,
		"customer_id": 68,
		"customer_first": "Willow",
		"customer_last": "Potter",
		"vendor_id": 68,
		"vendor_name": "Auctor Associates",
		"vendor_description": "magnis dis parturient montes, nascetur ridiculus mus. Donec",
		"vendor_contact_id": 68,
		"vendor_contact_first": "Elton",
		"vendor_contact_last": "Raymond"
	},
	{
		"id": 69,
		"created_time": "2021-04-12",
		"modified_time": "2020-04-03",
		"price": 2,
		"quantity": 6,
		"total": 7,
		"customer_id": 69,
		"customer_first": "Hyacinth",
		"customer_last": "Gaines",
		"vendor_id": 69,
		"vendor_name": "Risus Nulla Ltd",
		"vendor_description": "nisl sem, consequat nec, mollis vitae, posuere at,",
		"vendor_contact_id": 69,
		"vendor_contact_first": "Ishmael",
		"vendor_contact_last": "Carroll"
	},
	{
		"id": 70,
		"created_time": "2022-01-10",
		"modified_time": "2021-05-11",
		"price": 1,
		"quantity": 4,
		"total": 5,
		"customer_id": 70,
		"customer_first": "Faith",
		"customer_last": "Reese",
		"vendor_id": 70,
		"vendor_name": "Eleifend Nunc Risus Ltd",
		"vendor_description": "Integer urna. Vivamus molestie dapibus ligula.",
		"vendor_contact_id": 70,
		"vendor_contact_first": "Kai",
		"vendor_contact_last": "Bray"
	},
	{
		"id": 71,
		"created_time": "2020-07-31",
		"modified_time": "2021-02-11",
		"price": 3,
		"quantity": 6,
		"total": 3,
		"customer_id": 71,
		"customer_first": "Orlando",
		"customer_last": "Hale",
		"vendor_id": 71,
		"vendor_name": "Blandit Viverra Donec LLC",
		"vendor_description": "eget magna. Suspendisse tristique",
		"vendor_contact_id": 71,
		"vendor_contact_first": "Martina",
		"vendor_contact_last": "Whitaker"
	},
	{
		"id": 72,
		"created_time": "2020-10-23",
		"modified_time": "2022-01-03",
		"price": 8,
		"quantity": 6,
		"total": 8,
		"customer_id": 72,
		"customer_first": "Owen",
		"customer_last": "Wright",
		"vendor_id": 72,
		"vendor_name": "Vitae Purus Company",
		"vendor_description": "augue eu tellus. Phasellus elit",
		"vendor_contact_id": 72,
		"vendor_contact_first": "Calvin",
		"vendor_contact_last": "Alexander"
	},
	{
		"id": 73,
		"created_time": "2021-01-28",
		"modified_time": "2020-12-18",
		"price": 10,
		"quantity": 5,
		"total": 3,
		"customer_id": 73,
		"customer_first": "Colby",
		"customer_last": "Gould",
		"vendor_id": 73,
		"vendor_name": "Pellentesque Limited",
		"vendor_description": "volutpat nunc sit amet metus.",
		"vendor_contact_id": 73,
		"vendor_contact_first": "Callum",
		"vendor_contact_last": "Wiggins"
	},
	{
		"id": 74,
		"created_time": "2021-06-10",
		"modified_time": "2020-06-16",
		"price": 7,
		"quantity": 2,
		"total": 10,
		"customer_id": 74,
		"customer_first": "Wallace",
		"customer_last": "Mcgowan",
		"vendor_id": 74,
		"vendor_name": "Felis Nulla Tempor Incorporated",
		"vendor_description": "dolor. Fusce feugiat. Lorem ipsum dolor sit amet, consectetuer adipiscing",
		"vendor_contact_id": 74,
		"vendor_contact_first": "Venus",
		"vendor_contact_last": "Powers"
	},
	{
		"id": 75,
		"created_time": "2021-12-06",
		"modified_time": "2020-10-26",
		"price": 5,
		"quantity": 5,
		"total": 3,
		"customer_id": 75,
		"customer_first": "Lars",
		"customer_last": "Melendez",
		"vendor_id": 75,
		"vendor_name": "Sit Incorporated",
		"vendor_description": "lectus",
		"vendor_contact_id": 75,
		"vendor_contact_first": "Otto",
		"vendor_contact_last": "Ingram"
	},
	{
		"id": 76,
		"created_time": "2021-04-21",
		"modified_time": "2021-06-09",
		"price": 1,
		"quantity": 5,
		"total": 4,
		"customer_id": 76,
		"customer_first": "Byron",
		"customer_last": "Mcdowell",
		"vendor_id": 76,
		"vendor_name": "Et Consulting",
		"vendor_description": "tincidunt tempus risus. Donec egestas. Duis ac arcu. Nunc mauris.",
		"vendor_contact_id": 76,
		"vendor_contact_first": "Carly",
		"vendor_contact_last": "Chapman"
	},
	{
		"id": 77,
		"created_time": "2021-06-05",
		"modified_time": "2021-06-24",
		"price": 2,
		"quantity": 10,
		"total": 7,
		"customer_id": 77,
		"customer_first": "Carissa",
		"customer_last": "Sargent",
		"vendor_id": 77,
		"vendor_name": "Integer LLP",
		"vendor_description": "feugiat tellus lorem eu",
		"vendor_contact_id": 77,
		"vendor_contact_first": "Cade",
		"vendor_contact_last": "Horn"
	},
	{
		"id": 78,
		"created_time": "2020-11-01",
		"modified_time": "2021-07-27",
		"price": 6,
		"quantity": 7,
		"total": 6,
		"customer_id": 78,
		"customer_first": "Nathan",
		"customer_last": "Mcdaniel",
		"vendor_id": 78,
		"vendor_name": "Placerat Velit LLP",
		"vendor_description": "lobortis augue scelerisque mollis. Phasellus",
		"vendor_contact_id": 78,
		"vendor_contact_first": "Elijah",
		"vendor_contact_last": "Craft"
	},
	{
		"id": 79,
		"created_time": "2022-02-01",
		"modified_time": "2020-06-25",
		"price": 3,
		"quantity": 4,
		"total": 3,
		"customer_id": 79,
		"customer_first": "Stewart",
		"customer_last": "Gutierrez",
		"vendor_id": 79,
		"vendor_name": "Neque Consulting",
		"vendor_description": "dolor",
		"vendor_contact_id": 79,
		"vendor_contact_first": "Fitzgerald",
		"vendor_contact_last": "Bentley"
	},
	{
		"id": 80,
		"created_time": "2021-04-12",
		"modified_time": "2020-04-27",
		"price": 2,
		"quantity": 4,
		"total": 8,
		"customer_id": 80,
		"customer_first": "Valentine",
		"customer_last": "Raymond",
		"vendor_id": 80,
		"vendor_name": "Purus Sapien Gravida Corp.",
		"vendor_description": "vitae sodales",
		"vendor_contact_id": 80,
		"vendor_contact_first": "Travis",
		"vendor_contact_last": "Sanders"
	},
	{
		"id": 81,
		"created_time": "2021-01-17",
		"modified_time": "2021-07-04",
		"price": 7,
		"quantity": 6,
		"total": 3,
		"customer_id": 81,
		"customer_first": "Lillith",
		"customer_last": "Baldwin",
		"vendor_id": 81,
		"vendor_name": "Nascetur Ridiculus PC",
		"vendor_description": "at augue",
		"vendor_contact_id": 81,
		"vendor_contact_first": "Harding",
		"vendor_contact_last": "Kelley"
	},
	{
		"id": 82,
		"created_time": "2021-03-07",
		"modified_time": "2021-11-04",
		"price": 8,
		"quantity": 5,
		"total": 5,
		"customer_id": 82,
		"customer_first": "Kiona",
		"customer_last": "Harding",
		"vendor_id": 82,
		"vendor_name": "Diam Lorem Auctor Industries",
		"vendor_description": "quis massa. Mauris vestibulum, neque sed dictum eleifend, nunc risus",
		"vendor_contact_id": 82,
		"vendor_contact_first": "Maryam",
		"vendor_contact_last": "Jordan"
	},
	{
		"id": 83,
		"created_time": "2021-04-15",
		"modified_time": "2021-07-24",
		"price": 1,
		"quantity": 6,
		"total": 4,
		"customer_id": 83,
		"customer_first": "Kelsey",
		"customer_last": "Griffith",
		"vendor_id": 83,
		"vendor_name": "Nunc Institute",
		"vendor_description": "malesuada augue ut lacus. Nulla tincidunt, neque vitae semper",
		"vendor_contact_id": 83,
		"vendor_contact_first": "Allegra",
		"vendor_contact_last": "Wells"
	},
	{
		"id": 84,
		"created_time": "2020-04-12",
		"modified_time": "2021-05-29",
		"price": 7,
		"quantity": 10,
		"total": 8,
		"customer_id": 84,
		"customer_first": "Kane",
		"customer_last": "Christensen",
		"vendor_id": 84,
		"vendor_name": "Primis In Faucibus Incorporated",
		"vendor_description": "bibendum ullamcorper. Duis",
		"vendor_contact_id": 84,
		"vendor_contact_first": "Zeus",
		"vendor_contact_last": "Manning"
	},
	{
		"id": 85,
		"created_time": "2020-10-21",
		"modified_time": "2021-08-31",
		"price": 7,
		"quantity": 7,
		"total": 6,
		"customer_id": 85,
		"customer_first": "Hamilton",
		"customer_last": "Willis",
		"vendor_id": 85,
		"vendor_name": "Orci Quis Lectus Corporation",
		"vendor_description": "sem molestie sodales. Mauris",
		"vendor_contact_id": 85,
		"vendor_contact_first": "Zelenia",
		"vendor_contact_last": "Talley"
	},
	{
		"id": 86,
		"created_time": "2021-10-26",
		"modified_time": "2020-06-13",
		"price": 5,
		"quantity": 7,
		"total": 9,
		"customer_id": 86,
		"customer_first": "Madeson",
		"customer_last": "Albert",
		"vendor_id": 86,
		"vendor_name": "Aliquam Enim Foundation",
		"vendor_description": "egestas",
		"vendor_contact_id": 86,
		"vendor_contact_first": "Malcolm",
		"vendor_contact_last": "White"
	},
	{
		"id": 87,
		"created_time": "2020-12-09",
		"modified_time": "2020-08-13",
		"price": 3,
		"quantity": 7,
		"total": 1,
		"customer_id": 87,
		"customer_first": "Illana",
		"customer_last": "Mckenzie",
		"vendor_id": 87,
		"vendor_name": "Nunc Limited",
		"vendor_description": "elit. Curabitur sed tortor. Integer aliquam adipiscing lacus. Ut",
		"vendor_contact_id": 87,
		"vendor_contact_first": "Harper",
		"vendor_contact_last": "Mclaughlin"
	},
	{
		"id": 88,
		"created_time": "2020-12-12",
		"modified_time": "2021-01-31",
		"price": 7,
		"quantity": 3,
		"total": 2,
		"customer_id": 88,
		"customer_first": "Vernon",
		"customer_last": "Obrien",
		"vendor_id": 88,
		"vendor_name": "Placerat Eget Venenatis Corp.",
		"vendor_description": "Suspendisse aliquet molestie tellus.",
		"vendor_contact_id": 88,
		"vendor_contact_first": "Erasmus",
		"vendor_contact_last": "Tucker"
	},
	{
		"id": 89,
		"created_time": "2020-03-11",
		"modified_time": "2022-01-14",
		"price": 6,
		"quantity": 5,
		"total": 7,
		"customer_id": 89,
		"customer_first": "Ali",
		"customer_last": "Noble",
		"vendor_id": 89,
		"vendor_name": "Semper Rutrum Fusce Ltd",
		"vendor_description": "Donec non justo. Proin non massa non ante bibendum",
		"vendor_contact_id": 89,
		"vendor_contact_first": "Amy",
		"vendor_contact_last": "Roach"
	},
	{
		"id": 90,
		"created_time": "2021-07-29",
		"modified_time": "2020-10-06",
		"price": 3,
		"quantity": 10,
		"total": 2,
		"customer_id": 90,
		"customer_first": "Cyrus",
		"customer_last": "Gilmore",
		"vendor_id": 90,
		"vendor_name": "Duis Gravida LLP",
		"vendor_description": "at augue id ante",
		"vendor_contact_id": 90,
		"vendor_contact_first": "Rhea",
		"vendor_contact_last": "Huff"
	},
	{
		"id": 91,
		"created_time": "2021-01-08",
		"modified_time": "2021-09-27",
		"price": 3,
		"quantity": 9,
		"total": 5,
		"customer_id": 91,
		"customer_first": "Athena",
		"customer_last": "Wright",
		"vendor_id": 91,
		"vendor_name": "Purus Inc.",
		"vendor_description": "est. Mauris eu turpis. Nulla aliquet. Proin velit.",
		"vendor_contact_id": 91,
		"vendor_contact_first": "Stone",
		"vendor_contact_last": "English"
	},
	{
		"id": 92,
		"created_time": "2021-06-22",
		"modified_time": "2020-04-05",
		"price": 5,
		"quantity": 3,
		"total": 6,
		"customer_id": 92,
		"customer_first": "Lael",
		"customer_last": "Mullen",
		"vendor_id": 92,
		"vendor_name": "Egestas Hendrerit Neque LLC",
		"vendor_description": "ornare, libero at auctor ullamcorper, nisl",
		"vendor_contact_id": 92,
		"vendor_contact_first": "Hayden",
		"vendor_contact_last": "Hutchinson"
	},
	{
		"id": 93,
		"created_time": "2020-04-07",
		"modified_time": "2020-10-21",
		"price": 5,
		"quantity": 7,
		"total": 10,
		"customer_id": 93,
		"customer_first": "Allistair",
		"customer_last": "Stanley",
		"vendor_id": 93,
		"vendor_name": "Urna Et Arcu LLP",
		"vendor_description": "In ornare sagittis felis. Donec tempor,",
		"vendor_contact_id": 93,
		"vendor_contact_first": "Erich",
		"vendor_contact_last": "Perez"
	},
	{
		"id": 94,
		"created_time": "2022-01-19",
		"modified_time": "2021-11-02",
		"price": 6,
		"quantity": 3,
		"total": 3,
		"customer_id": 94,
		"customer_first": "Rogan",
		"customer_last": "Crawford",
		"vendor_id": 94,
		"vendor_name": "Integer Sem Elit Industries",
		"vendor_description": "eu",
		"vendor_contact_id": 94,
		"vendor_contact_first": "Caldwell",
		"vendor_contact_last": "Melton"
	},
	{
		"id": 95,
		"created_time": "2020-07-09",
		"modified_time": "2020-03-28",
		"price": 10,
		"quantity": 10,
		"total": 3,
		"customer_id": 95,
		"customer_first": "Chaim",
		"customer_last": "Patrick",
		"vendor_id": 95,
		"vendor_name": "Felis Adipiscing Institute",
		"vendor_description": "mattis semper, dui lectus rutrum",
		"vendor_contact_id": 95,
		"vendor_contact_first": "Ulla",
		"vendor_contact_last": "Bradley"
	},
	{
		"id": 96,
		"created_time": "2021-09-09",
		"modified_time": "2021-09-19",
		"price": 7,
		"quantity": 3,
		"total": 8,
		"customer_id": 96,
		"customer_first": "Caryn",
		"customer_last": "Moreno",
		"vendor_id": 96,
		"vendor_name": "Metus Facilisis Incorporated",
		"vendor_description": "erat eget ipsum. Suspendisse sagittis. Nullam vitae",
		"vendor_contact_id": 96,
		"vendor_contact_first": "Winifred",
		"vendor_contact_last": "Blanchard"
	},
	{
		"id": 97,
		"created_time": "2021-10-26",
		"modified_time": "2021-09-10",
		"price": 2,
		"quantity": 3,
		"total": 8,
		"customer_id": 97,
		"customer_first": "Timothy",
		"customer_last": "Calhoun",
		"vendor_id": 97,
		"vendor_name": "Semper Incorporated",
		"vendor_description": "posuere vulputate, lacus.",
		"vendor_contact_id": 97,
		"vendor_contact_first": "Pamela",
		"vendor_contact_last": "Coleman"
	},
	{
		"id": 98,
		"created_time": "2021-03-14",
		"modified_time": "2020-06-13",
		"price": 4,
		"quantity": 5,
		"total": 4,
		"customer_id": 98,
		"customer_first": "Naomi",
		"customer_last": "Trevino",
		"vendor_id": 98,
		"vendor_name": "Nec Urna Limited",
		"vendor_description": "in",
		"vendor_contact_id": 98,
		"vendor_contact_first": "Nicole",
		"vendor_contact_last": "Parker"
	},
	{
		"id": 99,
		"created_time": "2021-08-11",
		"modified_time": "2021-06-27",
		"price": 4,
		"quantity": 7,
		"total": 4,
		"customer_id": 99,
		"customer_first": "Chase",
		"customer_last": "Walter",
		"vendor_id": 99,
		"vendor_name": "Egestas Blandit Nam Associates",
		"vendor_description": "nibh. Aliquam ornare, libero at auctor ullamcorper, nisl arcu iaculis",
		"vendor_contact_id": 99,
		"vendor_contact_first": "Kay",
		"vendor_contact_last": "Hodges"
	},
	{
		"id": 100,
		"created_time": "2021-09-24",
		"modified_time": "2021-03-31",
		"price": 10,
		"quantity": 1,
		"total": 2,
		"customer_id": 100,
		"customer_first": "Keith",
		"customer_last": "Peck",
		"vendor_id": 100,
		"vendor_name": "Arcu Associates",
		"vendor_description": "sollicitudin orci sem eget massa. Suspendisse eleifend. Cras",
		"vendor_contact_id": 100,
		"vendor_contact_first": "Elmo",
		"vendor_contact_last": "Woodward"
	}
]
`

type MapperBenchmarkJsonRow struct {
	Id                 int    `json:"id"`
	CreatedTime        string `json:"created_time"`
	ModifiedTime       string `json:"modified_time"`
	Price              int    `json:"price"`
	Quantity           int    `json:"quantity"`
	Total              int    `json:"total"`
	CustomerId         int    `json:"customer_id"`
	CustomerFirst      string `json:"customer_first"`
	CustomerLast       string `json:"customer_last"`
	VendorId           int    `json:"vendor_id"`
	VendorName         string `json:"vendor_name"`
	VendorDescription  string `json:"vendor_description"`
	VendorContactId    int    `json:"vendor_contact_id"`
	VendorContactFirst string `json:"vendor_contact_first"`
	VendorContactLast  string `json:"vendor_contact_last"`
}

func loadBenchmarkMapperData(b *testing.B) ([]MapperBenchmarkJsonRow, int) {
	dest := []MapperBenchmarkJsonRow{}
	err := json.Unmarshal([]byte(benchmarkMapperJson), &dest)
	if err != nil {
		b.Fatalf("Unpacking JSON: %v", err.Error())
	}
	size := len(dest)
	return dest, size
}
func BenchmarkMapperBaseline(b *testing.B) {
	rows, size := loadBenchmarkMapperData(b)
	//
	type Person struct {
		Id    int
		First string
		Last  string
	}
	type Vendor struct {
		Id          int
		Name        string
		Description string
		Contact     Person
	}
	type T struct {
		Id           int
		CreatedTime  string
		ModifiedTime string
		Price        int
		Quantity     int
		Total        int
		Customer     Person
		Vendor       Vendor
	}
	//
	b.ResetTimer()
	//
	for k := 0; k < b.N; k++ {
		row := rows[k%size]
		dest := new(T)
		//
		dest.Id = row.Id
		dest.CreatedTime = row.CreatedTime
		dest.ModifiedTime = row.ModifiedTime
		dest.Price = row.Price
		dest.Quantity = row.Quantity
		dest.Total = row.Total
		//
		dest.Customer.Id = row.CustomerId
		dest.Customer.First = row.CustomerFirst
		dest.Customer.Last = row.CustomerLast
		//
		dest.Vendor.Id = row.VendorId
		dest.Vendor.Name = row.VendorName
		dest.Vendor.Description = row.VendorDescription
		dest.Vendor.Contact.Id = row.VendorContactId
		dest.Vendor.Contact.First = row.VendorContactFirst
		dest.Vendor.Contact.Last = row.VendorContactLast
	}
}

func BenchmarkMapperJsonUnmarshal(b *testing.B) {
	rowsDecoded, size := loadBenchmarkMapperData(b)
	var rows []string
	for _, decoded := range rowsDecoded {
		if encoded, err := json.Marshal(decoded); err != nil {
			b.Fatalf("During json.Marshal: %v", err.Error())
		} else {
			rows = append(rows, string(encoded))
		}

	}
	//
	type CustomerContact struct {
		Id    int    `json:"customer_id"`
		First string `json:"customer_first"`
		Last  string `json:"customer_last"`
	}
	type VendorContact struct {
		Id    int    `json:"vendor_contact_id"`
		First string `json:"vendor_contact_first"`
		Last  string `json:"vendor_contact_last"`
	}
	type Vendor struct {
		Id          int    `json:"vendor_id"`
		Name        string `json:"vendor_name"`
		Description string `json:"vendor_description"`
		VendorContact
	}
	type T struct {
		Id           int    `json:"id"`
		CreatedTime  string `json:"created_time"`
		ModifiedTime string `json:"modified_time"`
		Price        int    `json:"price"`
		Quantity     int    `json:"quantity"`
		Total        int    `json:"total"`

		CustomerContact
		Vendor
	}
	//
	b.ResetTimer()
	//
	for k := 0; k < b.N; k++ {
		row := rows[k%size]
		dest := new(T)
		if err := json.Unmarshal([]byte(row), dest); err != nil {
			b.Fatalf("During json.Unmarshal: %v", err.Error())
		}
	}
}
func BenchmarkMapperBoundMappingNoRebind(b *testing.B) {
	rows, size := loadBenchmarkMapperData(b)
	//
	type Common struct {
		Id int
	}
	type Timestamps struct {
		CreatedTime  string
		ModifiedTime string
	}
	type Person struct {
		Common
		Timestamps // Not used but present anyways
		First      string
		Last       string
	}
	type Vendor struct {
		Common
		Timestamps  // Not used but present anyways
		Name        string
		Description string
		Contact     Person
	}
	type T struct {
		Common
		Timestamps
		//
		Price    int
		Quantity int
		Total    int
		//
		Customer Person
		Vendor   Vendor
	}
	//
	b.ResetTimer()
	//
	mapper := &set.Mapper{
		Elevated: set.NewTypeList(Common{}, Timestamps{}),
		Join:     "_",
	}
	//
	for k := 0; k < b.N; k++ {
		row := rows[k%size]
		dest := new(T)
		bound := mapper.Bind(&dest)
		//
		bound.Set("Id", row.Id)
		bound.Set("CreatedTime", row.CreatedTime)
		bound.Set("ModifiedTime", row.ModifiedTime)
		bound.Set("Price", row.Price)
		bound.Set("Quantity", row.Quantity)
		bound.Set("Total", row.Total)
		//
		bound.Set("Customer_Id", row.CustomerId)
		bound.Set("Customer_First", row.CustomerFirst)
		bound.Set("Customer_Last", row.CustomerLast)
		//
		bound.Set("Vendor_Id", row.VendorId)
		bound.Set("Vendor_Name", row.VendorName)
		bound.Set("Vendor_Description", row.VendorDescription)
		bound.Set("Vendor_Contact_Id", row.VendorContactId)
		bound.Set("Vendor_Contact_First", row.VendorContactFirst)
		bound.Set("Vendor_Contact_Last", row.VendorContactLast)
		//
		if err := bound.Err(); err != nil {
			b.Fatalf("Unable to set: %v", err.Error())
		}
	}
}

func BenchmarkMapperBoundMappingWithRebind(b *testing.B) {
	rows, size := loadBenchmarkMapperData(b)
	//
	type Common struct {
		Id int
	}
	type Timestamps struct {
		CreatedTime  string
		ModifiedTime string
	}
	type Person struct {
		Common
		Timestamps // Not used but present anyways
		First      string
		Last       string
	}
	type Vendor struct {
		Common
		Timestamps  // Not used but present anyways
		Name        string
		Description string
		Contact     Person
	}
	type T struct {
		Common
		Timestamps
		//
		Price    int
		Quantity int
		Total    int
		//
		Customer Person
		Vendor   Vendor
	}
	//
	b.ResetTimer()
	//
	mapper := &set.Mapper{
		Elevated: set.NewTypeList(Common{}, Timestamps{}),
		Join:     "_",
	}
	//
	dest := new(T)
	bound := mapper.Bind(&dest)
	//
	for k := 0; k < b.N; k++ {
		row := rows[k%size]
		dest = new(T)
		bound.Rebind(&dest)
		//
		bound.Set("Id", row.Id)
		bound.Set("CreatedTime", row.CreatedTime)
		bound.Set("ModifiedTime", row.ModifiedTime)
		bound.Set("Price", row.Price)
		bound.Set("Quantity", row.Quantity)
		bound.Set("Total", row.Total)
		//
		bound.Set("Customer_Id", row.CustomerId)
		bound.Set("Customer_First", row.CustomerFirst)
		bound.Set("Customer_Last", row.CustomerLast)
		//
		bound.Set("Vendor_Id", row.VendorId)
		bound.Set("Vendor_Name", row.VendorName)
		bound.Set("Vendor_Description", row.VendorDescription)
		bound.Set("Vendor_Contact_Id", row.VendorContactId)
		bound.Set("Vendor_Contact_First", row.VendorContactFirst)
		bound.Set("Vendor_Contact_Last", row.VendorContactLast)
		//
		if err := bound.Err(); err != nil {
			b.Fatalf("Unable to set: %v", err.Error())
		}
	}
}

func BenchmarkMapperPreparedMapping(b *testing.B) {
	rows, size := loadBenchmarkMapperData(b)
	//
	type Common struct {
		Id int
	}
	type Timestamps struct {
		CreatedTime  string
		ModifiedTime string
	}
	type Person struct {
		Common
		Timestamps // Not used but present anyways
		First      string
		Last       string
	}
	type Vendor struct {
		Common
		Timestamps  // Not used but present anyways
		Name        string
		Description string
		Contact     Person
	}
	type T struct {
		Common
		Timestamps
		//
		Price    int
		Quantity int
		Total    int
		//
		Customer Person
		Vendor   Vendor
	}
	//
	b.ResetTimer()
	//
	mapper := &set.Mapper{
		Elevated: set.NewTypeList(Common{}, Timestamps{}),
		Join:     "_",
	}
	//
	dest := new(T)
	prepared := mapper.Prepare(&dest)
	err := prepared.Plan(
		"Id", "CreatedTime", "ModifiedTime",
		"Price", "Quantity", "Total",
		"Customer_Id", "Customer_First", "Customer_Last",
		"Vendor_Id", "Vendor_Name", "Vendor_Description", "Vendor_Contact_Id", "Vendor_Contact_First", "Vendor_Contact_Last")
	if err != nil {
		b.Fatalf("error preparing plan %v", err.Error())
	}
	//
	for k := 0; k < b.N; k++ {
		row := rows[k%size]
		dest = new(T)
		prepared.Rebind(&dest)
		//
		prepared.Set(row.Id)
		prepared.Set(row.CreatedTime)
		prepared.Set(row.ModifiedTime)
		prepared.Set(row.Price)
		prepared.Set(row.Quantity)
		prepared.Set(row.Total)
		//
		prepared.Set(row.CustomerId)
		prepared.Set(row.CustomerFirst)
		prepared.Set(row.CustomerLast)
		//
		prepared.Set(row.VendorId)
		prepared.Set(row.VendorName)
		prepared.Set(row.VendorDescription)
		prepared.Set(row.VendorContactId)
		prepared.Set(row.VendorContactFirst)
		prepared.Set(row.VendorContactLast)
		//
		if err := prepared.Err(); err != nil {
			b.Fatalf("Unable to set: %v", err.Error())
		}
	}
}

func BenchmarkValue(b *testing.B) { // TODO MOVE TO DIFFERENT FILE
	type Common struct {
		Id int
	}
	type Timestamps struct {
		CreatedTime  string
		ModifiedTime string
	}
	type Person struct {
		*Common
		*Timestamps // Not used but present anyways
		First       string
		Last        string
	}
	type Vendor struct {
		*Common
		*Timestamps // Not used but present anyways
		Name        string
		Description string
		Contact     Person
	}
	type T struct {
		*Common
		*Timestamps
		//
		Price    int
		Quantity int
		Total    int
		//
		Customer Person
		Vendor   Vendor
	}
	//
	for k := 0; k < b.N; k++ {
		dest := new(T)
		set.V(&dest)
	}
}

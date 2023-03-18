package main

import "testing"



func BenchmarkSolveChallenge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		chalValues := fetchCapChallenge(`{"do":["sid|855a6699-bb9a-11ed-bbff-626c70424d52","pnf|cu","cls|32060492575082180333|23765544575816579762","sts|1678050785554","wcs|cg2gboauqq6urpbp0480","drc|6412","cts|855a6a77-bb9a-11ed-bbff-626c70424d52|false","cs|99006a3a7a4e999df3c2e5de2b63611201318ad0dbbb050765d0c8a44bf201c3","vid|67932318-bb9a-11ed-89dd-436342496472|31536000|false","cp|1|35bc35ef77d4d787c8f08d309e88a5c3c2d0ef6cbc05abc4b5f0b097b8|d6af800afdeb6718b9207fd11ff6ea7718142f6d81ca79fd0d9f0da42168c031|25|false","ci|1|85628f20-bb9a-11ed-b230-9124f8607fdd|1056|ee0addfef8ccf455eddcf27b24c993887a5731ec83bb00b4d81e01f1f4e2314270372b7f80313bbe706fde46b9ab548dcfd596d6d8aee51352984821ed4b3a1a󠄻󠄺󠄸󠄸|0|NA","sff|cc|60|U2FtZVNpdGU9TGF4Ow==","sff|idp_c|60|1,s","sff|rf|60|1","sff|fp|60|1"]}`)

		solveChallenge(chalValues, 30)
	}
}

func BenchmarkSolveChallengeWithCoroutines(b *testing.B) {
	for i := 0; i < b.N; i++ {
		chalValues := fetchCapChallenge(`{"do":["sid|855a6699-bb9a-11ed-bbff-626c70424d52","pnf|cu","cls|32060492575082180333|23765544575816579762","sts|1678050785554","wcs|cg2gboauqq6urpbp0480","drc|6412","cts|855a6a77-bb9a-11ed-bbff-626c70424d52|false","cs|99006a3a7a4e999df3c2e5de2b63611201318ad0dbbb050765d0c8a44bf201c3","vid|67932318-bb9a-11ed-89dd-436342496472|31536000|false","cp|1|35bc35ef77d4d787c8f08d309e88a5c3c2d0ef6cbc05abc4b5f0b097b8|d6af800afdeb6718b9207fd11ff6ea7718142f6d81ca79fd0d9f0da42168c031|25|false","ci|1|85628f20-bb9a-11ed-b230-9124f8607fdd|1056|ee0addfef8ccf455eddcf27b24c993887a5731ec83bb00b4d81e01f1f4e2314270372b7f80313bbe706fde46b9ab548dcfd596d6d8aee51352984821ed4b3a1a󠄻󠄺󠄸󠄸|0|NA","sff|cc|60|U2FtZVNpdGU9TGF4Ow==","sff|idp_c|60|1,s","sff|rf|60|1","sff|fp|60|1"]}`)

		solveChallengeWithCoroutines(chalValues, 30)
	}
}

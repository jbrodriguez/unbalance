import test from 'ava'

import { humanBytes } from '../src/lib/utils'

test('Check humanBytes', t => {
	const expected = '1.00 TB'
	const actual = humanBytes(999716474880)

	t.deepEqual(actual, expected)
})

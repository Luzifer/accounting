export function formatNumber(number, thousandSep = ' ', decimalSep = '.', places = 2) {
  if (isNaN(number)) {
    return number
  }

  // Fix x.99999999999 database fuckups
  number = Math.round(number * Math.pow(10, 2)) / Math.pow(10, 2)

  let result = number < 0 ? '-' : ''
  number = Math.abs(number)
  if (number >= Number.MAX_SAFE_INTEGER) {
    return result + number.toFixed(places)
  }

  let place = Math.ceil(Math.log10(number))

  if (place < 3) {
    return result + number.toFixed(places).replace('.', decimalSep)
  }

  while (place--) {
    result += number / 10 ** place % 10 | 0
    if (place > 0 && place % 3 === 0) {
      result += thousandSep
    }
  }

  return result + decimalSep + number.toFixed(places).split('.')[1]
}

/**
 * Common code to derive a class from a numeric value
 *
 * @param {Number} num The value to choose the class from
 * @param {Array} extraClasses Extra classes to add to the output string
 * @param {String | null} positiveClass Class to use on positive numbers
 * @returns {String} Space separated combined class list
 */
export function classFromNumber(num, extraClasses = [], positiveClass = null) {
  const classes = extraClasses || []
  if (num < 0) {
    classes.push('text-danger')
  } else if (num === 0) {
    classes.push('text-muted')
  } else if (positiveClass) {
    classes.push(positiveClass)
  }

  return classes.join(' ')
}

/**
 * Parses the response to JSON and throws an exception in case the
 * request was non-2xx
 *
 * @param {Response} resp Response from a `fetch` request
 */
export function responseToJSON(resp) {
  if (resp.status > 299) {
    throw new Error(`non-2xx status code: ${resp.status}`)
  }

  if (resp.status === 204) {
    return null
  }

  return resp.json()
}

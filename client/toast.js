/** @format */

import Toastify from 'toastify-js'

function toast(type, html, duration) {
  Toastify({
    text: html,
    duration: 7000,
    close: duration > 10000,
    gravity: 'top',
    positionLeft: false,
    className: `toast-${type}`
  }).showToast()
}

export const info = (html, d = 7000) => toast('info', html, d)
export const warning = (html, d = 7000) => toast('warning', html, d)
export const error = (html, d = 7000) => toast('error', html, d)
export const success = (html, d = 7000) => toast('success', html, d)

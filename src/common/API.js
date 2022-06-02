import axios from 'axios';

export class APIResponse {
  constructor(result) {
    if (typeof (result) === 'string' || result instanceof String) {
      // custom error message
      this.message = result;

    } else {
      this.result = result;

      if (result.response !== undefined) {
        if (result.response.status !== undefined) {
          this.code = result.response.status;
        }
      }
      if (result.status !== undefined) {
        this.code = result.status;
      }
      if (result.data !== undefined) {
        this.payload = result.data;
      }
      if (this.message === undefined && result.message !== undefined) {
        this.message = result.message;
      }
    }
  }

  isSuccess() {
    return this.result !== undefined && this.code === 200;
  }

  isError() {
    return !this.isSuccess();
  }

  isUnauthorized() {
    return this.result !== undefined && this.code === 403;
  }

  isForbidden() {
    return this.result !== undefined && this.code === 401;
  }

  msg() {
    if (this.isSuccess()) {
      return 'OK';
    } else if (this.isUnauthorized()) {
      return 'Unauthorized';
    } else {
      if (this.code !== undefined) {
        return 'Server returned code ' + this.code;
      }
      if (this.message !== undefined) {
        return this.message;
      }
      return 'An error occured';
    }
  }
}

function handle(promise) {
  return promise
    .then((response) => {
      return new APIResponse(response);
    })
    .catch((err) => {
      return new APIResponse(err);
    })
}
function get(url) {
  return handle(axios.get(url));
}
function post(url, body) {
  return handle(axios.post(url, body));
}
function put(url, body) {
  return handle(axios.put(url, body));
}
function del(url) {
  return handle(axios.delete(url));
}

export function getDBInfo() {
  // obtain new admin token using a login token
  return get("/api/dbinfo");
}

export function getGroups() {
  // refresh admin token using a cookie
  return get("/api/groups");
}

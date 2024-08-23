import request from "@/utils/request";
export function getSubs(){
  return request({
    url: "/api/v1/subscription/get",
    method: "get",
  });
}

export function AddSub(data: any){
  return request({
    url: "/api/v1/subscription/add",
    method: "post",
    data,
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
}
export function DelSub(data: any){
  return request({
    url: "/api/v1/subscription/delete",
    method: "delete",
    params: data,
  });
}

export function UpdateSub(data: any){
  return request({
    url: "/api/v1/subscription/update",
    method: "post",
    data,
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
}

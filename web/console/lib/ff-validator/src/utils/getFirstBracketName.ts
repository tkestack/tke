/**
 * 获取bracket的名称
 * @return string
 */
export const getBracketName = (keyName: string) => {
  let bracketIndex = keyName.indexOf('[');
  return bracketIndex > -1 ? keyName.slice(0, bracketIndex) : keyName;
};

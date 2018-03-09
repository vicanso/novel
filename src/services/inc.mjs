import genService from './gen';

const incService = genService('Inc');

async function inc(category) {
  const doc = await incService
    .findOneAndUpdate(
      {
        category,
      },
      {
        category,
        $inc: {
          value: 1,
        },
      },
      {
        new: true,
        upsert: true,
      },
    )
    .lean();
  return doc.value;
}

export async function getBookNo() {
  const no = await inc('bookNo');
  return no;
}

export async function getUserNo() {
  const no = await inc('userNo');
  return no;
}

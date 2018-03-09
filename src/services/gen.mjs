import _ from 'lodash';
import * as mongo from '../helpers/mongo';

/**
 * 往collection中插入一条记录
 *
 * @param {String} collection mongodb collection
 * @param {Object} data 要插入的记录
 */
function add(collection, data) {
  const Model = mongo.get(collection);
  return new Model(data).save().then(doc => doc.toJSON());
}

/**
 * 更新数据库记录(只更新一条), 参数与findOneAndUpdate一致
 *
 * @param {String} collection mongodb collection
 * @param {any} args
 */
function findOneAndUpdate(collection, ...args) {
  const Model = mongo.get(collection);
  return Model.findOneAndUpdate(...args);
}

/**
 * 查找数据，参数与mongoose find一致
 *
 * @param {String} collection mongodb collection
 * @param {any} args
 */
function find(collection, ...args) {
  const Model = mongo.get(collection);
  return Model.find(...args);
}

/**
 * 通过ID更新数据库记录
 *
 * @param {String} collection mongodb collection
 * @param {ObjectId} id  mongodb id
 * @param {Object} data 要更新的数据
 */
function findByIdAndUpdate(collection, id, data) {
  return findOneAndUpdate(
    collection,
    {
      _id: id,
    },
    data,
  );
}

/**
 * 查找一条记录
 *
 * @param {String} collection mongodb collection
 * @param {any} args
 */
function findOne(collection, ...args) {
  const Model = mongo.get(collection);
  return Model.findOne(...args);
}

/**
 * 通过id查找记录
 *
 * @param {String} collection mongodb collection
 * @param {ObjectId} id
 */
function findById(collection, id) {
  return findOne(collection, {
    _id: id,
  });
}

/**
 * Count
 *
 * @param {String} collection mongodb collection
 * @param {any} args
 */
function count(collection, ...args) {
  const Model = mongo.get(collection);
  return Model.count(...args);
}

export default function gen(collection) {
  const fns = {
    add,
    findOneAndUpdate,
    find,
    findByIdAndUpdate,
    findOne,
    findById,
    count,
  };
  const wrapper = {};
  _.forEach(fns, (fn, name) => {
    wrapper[name] = (...args) => fn(collection, ...args);
  });
  wrapper.Model = mongo.get(collection);
  return wrapper;
}
